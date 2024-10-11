import base64
from datetime import datetime, timezone
from .crypto import HashField, Hashable, Signer

def castDatetime(obj):
    if obj is None:
        return None
    if isinstance(obj, datetime):
        dt = obj.replace(microsecond=0)
    elif isinstance(obj, str):
        dt = castDatetime(datetime.fromisoformat(obj))
    else:
        raise RuntimeError(f"{obj} cannot be converted to a datetime object")
    if dt.tzname() is None:
       raise RuntimeError(f'Invalid datetime, no timezone "{obj}"')
    return dt

class User(Hashable):
    Username = HashField(i='n')
    Scope = HashField(i='s')
    Email = HashField(i='e')
    Key = HashField(i='k')

    def __init__(self, Username, Email, Key, Scope = None):
        self.Username = Username
        self.Email = Email
        self.Key = Key
        self.Scope = Scope

class Group(Hashable):
    Name = HashField(i='n')
    Scope = HashField(i='s')

    def __init__(self, Name, Scope=None):
        self.Name = Name
        self.Scope = Scope

class Module(Hashable):
    Name = HashField(i='n')
    Audience = HashField(i='a')
    Before = HashField(i='b')
    Reveal = HashField(i='r')

    def __init__(self, Name, Audience, Before, Reveal):
        self.Name = Name
        self.Audience = int(Audience)
        self.Before = castDatetime(Before)
        self.Reveal = castDatetime(Reveal if Reveal is not None else self.Before)

class Question(Hashable):
    Name = HashField(i='n')
    Before = HashField(i='b')
    Reveal = HashField(i='r')
    Grader = HashField(i='g')
    MaxScore = HashField(i='h')
    MinScore = HashField(i='m')
    MaxTry = HashField(i='t')

    def __init__(self, Name, Grader, MinScore, MaxScore, MaxTry, Before, Reveal):
        self.Name = Name
        self.Before = castDatetime(Before)
        self.Reveal = castDatetime(Reveal if Reveal is not None else self.Before)
        self.Grader = Grader
        self.MinScore = MinScore
        self.MaxScore = MaxScore
        self.MaxTry = MaxTry

class Submission(Hashable):
    Group    = HashField(i='g')
    Question = HashField(i='q')
    Data     = HashField(i='d')

    def __init__(self, Question, Group, Data):
        self.Question = Question
        self.Group = Group
        self.Data = base64.b64encode(Data).decode('utf-8')

