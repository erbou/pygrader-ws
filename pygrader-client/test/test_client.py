import pygrader.client as client
import json
import logging

conn = client.Client('https', 'localhost', 9443, key_pem_path='./client.sign.key.pem', cert_ca_path='./ca-chain.cert.pem', cert_pem_path = './client.bundle.pem')

def Result(s, m, data):
    Id = None
    print(f'"{data}"')
    if data is not None:
        data = json.loads(data)
        if isinstance(data, dict):
            Id = data.get('Id', None)
    print(16*'-')
    print(f'status:{s} msg:{m} -- {data}')
    return Id

try:
    print(64*'-' + '[list_user]' + 64*'-')
    s, m, data = conn.list_user(Email=conn.user)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_user]' + 64*'-')
    s, m, data = conn.create_user('Eric', conn.user)
    userId = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_user(Email)]' + 64*'-')
    s, m, data = conn.list_user(Email=conn.user)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_user]' + 64*'-')
    s, m, data = conn.list_user()
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[get_user]' + 64*'-')
    s, m, data = conn.get_user(userId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[update_user]' + 64*'-')
    s, m, data = conn.update_user(userId, Username='Emily', Email='emily@home.ch')
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_group]' + 64*'-')
    s, m, data = conn.create_group('Bouillet', 'owner')
    groupId = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_group]' + 64*'-')
    s, m, data = conn.list_group()
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[get_group]' + 64*'-')
    s, m, data = conn.get_group(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[update_group]' + 64*'-')
    s, m, data = conn.update_group(groupId, Name='Bouillets')
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_group(Name)]' + 64*'-')
    s, m, data = conn.list_group(Name='Bouillets')
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[group_add_user]' + 64*'-')
    s, m, data = conn.group_add_user(groupId, userId)
    userGroupId = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_module]' + 64*'-')
    s, m, data = conn.create_module('ENG209.Quiz.1')
    moduleId = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_module]' + 64*'-')
    s, m, data = conn.list_module()
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_question]' + 64*'-')
    s, m, data = conn.create_question(moduleId, 'q1', '2024-12-31T00:00:00+01:00', 'module1-grader', 6, 0)
    questionId = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_user]' + 64*'-')
    s, m, data = conn.delete_user(userId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_group]' + 64*'-')
    s, m, data = conn.delete_group(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_module]' + 64*'-')
    s, m, data = conn.delete_module(moduleId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_question]' + 64*'-')
    s, m, data = conn.delete_question(questionId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
