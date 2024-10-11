#!/usr/bin/python3

import http.client
import base64
import os
import re
import ssl
import json
import logging

from datetime import datetime, timezone
from cryptography import x509
import pygrader.crypto as crypto
import pygrader.models as models

class Client:
    """
        Create a new client.
        * Scheme must be http or https
        * Host, port are the grader's server name or IP and port number
        * Key_pem_path is the path to the private key used to sign the payloads of post and put requests
        * Key_password is the password of the private key used to sign the payloads (if PEM encrypted).
        * Cert_ca: is the path to the certificate chain used to authenticate the server
        * Cert_pem_path: is the path to the client certificate
        * Cert_key_pem_path: is the path to the private key of the client certificate (None, if the key is bundled with the certificate)
        * Cert_key_password: is the password used to encrypt the client certificate's key (if PEM encrypted).
        * Note: if key_pem_path is null, the signer will sign with the private key from the certificate.
    """
    def __init__(self,
            scheme, host, port,
            key_pem_path, key_password = None,
            cert_ca_path = None,
            cert_pem_path=None, cert_key_pem_path=None, cert_key_password = None,
            logLevel = logging.INFO,
        ):
        if key_pem_path is None:
            key_pem_path = cert_key_pem_path
            key_password = cert_key_password
        if key_pem_path is None:
            key_pem_path = cert_pem_path
        self.signer = crypto.Signer(key_pem_path, key_password)
        self.signer.key_fingerprint_len = 20
        self.host = host
        self.port = port
        self.scheme = scheme.lower()
        self.withDigest = False
        self.context = None
        self.client_certificate(cert_pem_path, cert_key_pem_path, cert_key_password)
        self.server_certificate(cert_ca_path)
        self._logger = logging.getLogger(__name__)
        self._logger.setLevel(logLevel)
        handler = logging.StreamHandler()
        handler.setLevel(logging.INFO)
        formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
        handler.setFormatter(formatter)
        self._logger.addHandler(handler)

    def client_certificate(self, cert_pem_path, cert_key_pem_path, cert_key_password):
        if cert_pem_path is None:
            return
        self.cert_pem = os.path.realpath(cert_pem_path)
        if cert_key_pem_path is not None:
            cert_key_pem_path = os.path.realpath(cert_key_pem_path)
        self.password = cert_key_password
        if self.context is None:
            self.context = ssl.create_default_context(ssl.Purpose.SERVER_AUTH)
        self.context.load_cert_chain(cert_pem_path, cert_key_pem_path, cert_key_password)
        self.cert_key_pem_path = cert_key_pem_path
        self.cert_pem_path = cert_pem_path
        with open(self.cert_pem_path, 'rb') as pem:
            cert = x509.load_pem_x509_certificate(pem.read())
        current_time = datetime.now(timezone.utc)
        if cert.not_valid_before_utc > current_time or cert.not_valid_after_utc < current_time:
            raise(f'Certificate has expired {cert.not_valid_before_utc} < {current_time} < {cert.not_valid_after_utc}')
        self._user = cert.subject.get_attributes_for_oid(x509.NameOID.COMMON_NAME)[0].value

    def server_certificate(self, cert_ca_path):
        if cert_ca_path is None:
            return
        cert_ca_path = os.path.realpath(cert_ca_path)
        if self.context is None:
            self.context = ssl.create_default_context(ssl.Purpose.SERVER_AUTH)
        self.ca_path = cert_ca_path
        self.context.load_default_certs()
        self.context.load_verify_locations(cert_ca_path)
        self.context.verify_mode = ssl.CERT_REQUIRED
        self.context.check_hostname = True

    @property
    def user(self):
        return self._user

    @property
    def logger(self):
        return self._logger

    def buildUrl(self, path, query=[]):
        query = '&'.join([ f'{p[0]}={p[1]}' for p in query if p[1] is not None])
        return path.rstrip('/') + ('?' if len(query)>0 else '/') + query

    def commit(self, Method, Path, Body=None):
        self.logger.info(f'{Method.upper()} {Path} {Body}')
        try:
            if self.scheme == 'http':
                conn = http.client.HTTPConnection(self.host, port=self.port, timeout=10)
            else:
                conn = http.client.HTTPSConnection(self.host, port=self.port, timeout=10, context = self.context)
        except Exception as e:
            self.logger.error(f'Error creating connection {e}')
            return

        try:
            conn.request(Method.upper(),
                 Path,
                 body=Body,
                 headers={
                     "Content-Type": "application/json"
                 })
            response = conn.getresponse()
            data = response.read()
            return response.status, response.reason, data.decode()
        except Exception as e:
            self.logger.error(f'Error submitting request: {e}')
        finally:
            conn.close()
        return 500, None, None

    def create_user(self, Username, Email, Key=None, Scope=None):
        if Key is None:
            Key = self.signer.public_key_b64.decode()
        body = self.signer.signObject(models.User(Username, Email, Key, Scope), withDigest = self.withDigest)
        return self.commit('POST', f'/v1/user/', body)

    def get_user(self, Id):
        return self.commit('GET', f'/v1/user/{Id}')

    def update_user(self, Id, Username=None, Email=None, Key=None):
        body = self.signer.signObject(models.User(Username, Email, Key), params=[Id], withDigest = self.withDigest)
        return self.commit('PUT', f'/v1/user/{Id}', body)

    def list_user(self, Username=None, Email=None, Kid=None):
        url = self.buildUrl('/v1/user', [('username',Username),('email',Email),('kid',Kid)])
        return self.commit('GET', url)

    def delete_user(self, Id):
        return self.commit('DELETE', f'/v1/user/{Id}')

    def create_group(self, Name, Scope=None):
        body = self.signer.signObject(models.Group(Name, Scope), withDigest = self.withDigest)
        return self.commit('POST', f'/v1/group/', body)

    def get_group(self, Id):
        return self.commit('GET', f'/v1/group/{Id}')

    def update_group(self, Id, Name=None, Scope=None):
        body = self.signer.signObject(models.Group(Name, Scope), params=[Id], withDigest = self.withDigest)
        return self.commit('PUT', f'/v1/group/{Id}', body)

    def group_add_user(self, GroupId, UserId, Token=None):
        if Token is not None:
            return self.commit('POST', f'/v1/group/{GroupId}/user/{UserId}?secret={Token}')
        else:
            return self.commit('POST', f'/v1/group/{GroupId}/user/{UserId}')

    def group_list_user(self, GroupId):
        url = self.buildUrl(f'/v1/group/{GroupId}/user', [])
        return self.commit('GET', url)

    def group_remove_user(self, GroupId, UserId, Token=None):
        if Token is not None:
            return self.commit('DELETE', f'/v1/group/{GroupId}/user/{UserId}?secret={Token}')
        else:
            return self.commit('DELETE', f'/v1/group/{GroupId}/user/{UserId}')

    def group_add_subgroup(self, GroupId, SubgroupId, Token=None):
        if Token is not None:
            return self.commit('POST', f'/v1/group/{GroupId}/group/{SubgroupId}?secret={Token}')
        else:
            return self.commit('POST', f'/v1/group/{GroupId}/group/{SubgroupId}')

    def group_remove_subgroup(self, GroupId, SubgroupId, Token=None):
        if Token is not None:
            return self.commit('DELETE', f'/v1/group/{GroupId}/group/{SubgroupId}?secret={Token}')
        else:
            return self.commit('DELETE', f'/v1/group/{GroupId}/group/{SubgroupId}')

    def group_list_subgroup(self, GroupId):
        url = self.buildUrl(f'/v1/group/{GroupId}/sub', [])
        return self.commit('GET', url)

    def group_list_supergroup(self, GroupId):
        url = self.buildUrl(f'/v1/group/{GroupId}/sup', [])
        return self.commit('GET', url)

    def delete_group(self, Id):
        return self.commit('DELETE', f'/v1/group/{Id}')

    def list_group(self, Name=None):
        url = self.buildUrl('/v1/group', [('name',Name)])
        return self.commit('GET', url)

    def create_module(self, Name, Audience, Before, Reveal):
        body = self.signer.signObject(models.Module(Name, Audience, Before, Reveal), withDigest = self.withDigest)
        return self.commit('POST', f'/v1/module/', body)

    def get_module(self, Id):
        return self.commit('GET', f'/v1/module/{Id}')

    def update_module(self, Id, Name=None, Audience=None, Before=None, Reveal=None):
        body = self.signer.signObject(models.Module(Name, Audience, Before, Reveal), params=[Id], withDigest = self.withDigest)
        return self.commit('PUT', f'/v1/module/{Id}', body)

    def delete_module(self, Id):
        return self.commit('DELETE', f'/v1/module/{Id}')

    def list_module(self, Name=None):
        url = self.buildUrl('/v1/module', [('name',Name)])
        return self.commit('GET', url)

    def list_module_questions(self, Id):
        return self.commit('GET', f'/v1/module/{Id}/question')

    def create_question(self, ModuleId, Name, Grader, MinScore, MaxScore, MaxTry, Before, Reveal):
        body = self.signer.signObject(models.Question(Name, Grader, MinScore, MaxScore, MaxTry, Before, Reveal), params=[ModuleId], withDigest = self.withDigest)
        #return self.commit('POST', f'/v1/question/module/{ModuleId}', body)
        return self.commit('POST', f'/v1/module/{ModuleId}/question', body)

    def get_question(self, Id):
        return self.commit('GET', f'/v1/question/{Id}')

    def update_question(self, Id, Name=None, Grader=None, MaxScore=None, MinScore=None, MaxTry=None, Before=None, Reveal=None):
        body = self.signer.signObject(models.Question(Name, Grader, MaxScore, MinScore, MaxTry, Before, Reveal), params=[Id], withDigest = self.withDigest)
        return self.commit('PUT', f'/v1/question/{Id}', body)

    def delete_question(self, Id):
        return self.commit('DELETE', f'/v1/question/{Id}')

    def answer_question(self, Id, Group, Data):
        body = self.signer.signObject(models.Submission(Id, Group, Data), withDigest = self.withDigest)
        return self.commit('POST', f'/v1/question/{Id}/answer/', body)

