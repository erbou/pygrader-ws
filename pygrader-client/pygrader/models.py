import http.client
from datetime import datetime
from .crypto import HashField, Hashable, Signer

class User(Hashable):
    Username = HashField(i='n')
    Email = HashField(i='e')
    Key = HashField(i='k')

    def __init__(self, Username, Email, Key):
        self.Username = Username
        self.Email = Email
        self.Key = Key

class Group(Hashable):
    Name = HashField(i='n')
    Scope = HashField(i='s')
    Token = HashField(i='t')

    def __init__(self, Name, Scope=None, Token=None):
        self.Name = Name
        self.Scope = Scope
        self.Token = Token

class Module(Hashable):
    Name = HashField(i='n')

    def __init__(self, Name):
        self.Name = Name

class Question(Hashable):
    Name = HashField(i='n')
    Before = HashField(i='b')
    Grader = HashField(i='g')
    MaxScore = HashField(i='h')
    MinScore = HashField(i='m')

    def __init__(self, Name, Before, Grader, MaxScore, MinScore):
        if isinstance(Before, str):
            Before = datetime.fromisoformat(Before)
            if Before.tzname() is None:
                raise RuntimeError(f'Invalid datetime, no timezone "{Before}"')
        elif not isinstance(Before, datetime):
            raise RuntimeError(f'Invalid datetime type: {type(Before)}')
        self.Name = Name
        self.Before = Before
        self.Grader = Grader
        self.MinScore = MinScore
        self.MaxScore = MaxScore
