import pygrader.client as client
import json
import pickle
import logging
from datetime import datetime, timedelta

conn = client.Client('https', 'localhost', 9443,
	key_pem_path='./client.sign.key.pem', # Signing key. Use key from client bundle if None
	cert_ca_path='./ca-chain.cert.pem',   # CA path for server auth
	cert_pem_path = './client.bundle.pem' # Client key and x509 certificate for client auth
)

def Result(s, m, data):
    try:
        Body = json.loads(data)
        if isinstance(Body, dict):
           Id = Body.get('Id', None)
        else:
           Id = None
    except:
        print(f"Exception in {data}")
        Id = None
        Body = None
		
    print(16*'-')
    print(f'status:{s} msg:{m}  {data}')
    return Id, Body

try:
    print(64*'-' + '[list_user]' + 64*'-')
    s, m, data = conn.list_user(Email=conn.user)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_user]' + 64*'-')
    s, m, data = conn.create_user('John', conn.user)
    userId, _ = Result(s, m, data)
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
    s, m, data = conn.update_user(userId, Username='Alice', Email='alice@home.ch')
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
    print(64*'-' + '[create_group]' + 64*'-')
    s, m, data = conn.create_group('Smith', 'owner')
    groupId, Body = Result(s, m, data)
    groupSecret = Body['Token']
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
    s, m, data = conn.update_group(groupId, Name='Smiths')
    _, Body = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[group_add_user]' + 64*'-')
    s, m, data = conn.group_add_user(groupId, userId, groupSecret)
    userGroupId, _ = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[create_group]' + 64*'-')
    s, m, data = conn.create_group('Jones', 'owner')
    groupId2, Body = Result(s, m, data)
    groupSecret2 = Body['Token']
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[add_sub_group]' + 64*'-')
    s, m, data = conn.group_add_subgroup(groupId, groupId2, groupSecret)
    subGroupId, Body = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_group(Name)]' + 64*'-')
    s, m, data = conn.list_group(Name='Smiths')
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[list_subgroup()]' + 64*'-')
    s, m, data = conn.group_list_subgroup(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[list_supergroup()]' + 64*'-')
    s, m, data = conn.group_list_supergroup(groupId2)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + f'[get_group({groupId})]' + 64*'-')
    s, m, data = conn.get_group(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + f'[get_group({groupId2})]' + 64*'-')
    s, m, data = conn.get_group(groupId2)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[create_module]' + 64*'-')
    s, m, data = conn.create_module(Name='ENG209.Quiz.1', Audience=groupId, Before='2024-12-31T00:00:00+01:00', Reveal=None)
    moduleId, _ = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[get_module]' + 64*'-')
    s, m, data = conn.get_module(moduleId)
    Result(s, m, data)
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
    s, m, data = conn.create_question(moduleId, 'q1', Before='2024-12-31T00:00:00+01:00', Reveal=None, Grader='module1-grader', MaxScore=6, MinScore=0, MaxTry=2)
    questionId, _ = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[answer_question]' + 64*'-')
    s, m, data = conn.answer_question(questionId, groupId2, pickle.dumps({ "TEST": 1 }))
    Result(s, m, data)
    s, m, data = conn.answer_question(questionId, groupId2, pickle.dumps({ "TEST": 1 }))
    Result(s, m, data)
    s, m, data = conn.answer_question(questionId, groupId2, pickle.dumps({ "TEST": 2 }))
    Result(s, m, data)
    s, m, data = conn.answer_question(questionId, groupId2, pickle.dumps({ "TEST": 3 }))
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_user]' + 64*'-')
    s, m, data = conn.delete_user(userId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[remove_sub_group]' + 64*'-')
    s, m, data = conn.group_remove_subgroup(groupId, groupId2)
    subGroupId, Body = Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[list_subgroup()]' + 64*'-')
    s, m, data = conn.group_list_subgroup(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[list_supergroup()]' + 64*'-')
    s, m, data = conn.group_list_supergroup(groupId2)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + f'[get_group({groupId})]' + 64*'-')
    s, m, data = conn.get_group(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + f'[get_group({groupId2})]' + 64*'-')
    s, m, data = conn.get_group(groupId2)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')
    raise(e)

try:
    print(64*'-' + '[delete_group]' + 64*'-')
    s, m, data = conn.delete_group(groupId)
    Result(s, m, data)
except Exception as e:
    print(f'Failed {e}')

try:
    print(64*'-' + '[delete_group]' + 64*'-')
    s, m, data = conn.delete_group(groupId2)
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

