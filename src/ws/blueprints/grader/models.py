from __future__ import annotations
import datetime
from ws import db, DillPickleType
from sqlalchemy import ForeignKey, LargeBinary, Index, UniqueConstraint, Column, Table
from sqlalchemy.orm import Mapped, mapped_column, relationship

class _Base(db.Model):
    __abstract__ = True
    id : Mapped[int] = mapped_column(primary_key=True)
    created: Mapped[datetime.datetime] = mapped_column(default=datetime.datetime.utcnow)
    updated: Mapped[datetime.datetime] = mapped_column(default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow)

    @classmethod
    def get(cls,id:int):
        return db.session.get(cls,id)

association_group_user = Table(
    'association_group_user',
    _Base.metadata,
    Column('group', ForeignKey('group.id'), primary_key=True),
    Column('user', ForeignKey('user.id'), primary_key=True),
)

class AssociationModuleUser(_Base):
    module_id: Mapped[int] = mapped_column(ForeignKey('module.id'), primary_key=True)
    user_id: Mapped[int] = mapped_column(ForeignKey('user.id'), primary_key=True)
    group_id: Mapped[int] = mapped_column(ForeignKey('group.id'))
    module: Mapped['Module'] = relationship(back_populates='user_association')
    user: Mapped['User'] = relationship(back_populates='module_association')

class User(_Base):
    username: Mapped[str] = mapped_column(index=True, unique=True)
    email: Mapped[str] = mapped_column(unique=True)
    public_key: Mapped[bytes] = mapped_column(LargeBinary, unique=True)
    groups: Mapped[list[Group]] = relationship(secondary=association_group_user, back_populates='users')
    module_association: Mapped[list[AssociationModuleUser]] = relationship(back_populates='user')

    @classmethod
    def find(cls, name:str):
        return db.session.query(cls).filter_by(username=name)

    def __repr__(self):
        return f'id={self.id}, username={self.username}, key={self.public_key}, created={self.created}, updated={self.updated}'

class Group(_Base):
    name: Mapped[str] = mapped_column(index=True, unique=True)
    scope: Mapped[str] = mapped_column()
    answers: Mapped[list['Answer']] = relationship(back_populates='group', cascade='all, delete-orphan', passive_deletes=True)
    managed_modules: Mapped[list['Module']] = relationship(back_populates='admin', cascade='all, delete-orphan', passive_deletes=True)
    users: Mapped[list[User]] = relationship(secondary=association_group_user, back_populates='groups')

    @classmethod
    def find(cls, name:str):
        return db.session.query(cls).filter_by(name=name)

    def __repr__(self):
        return f'id={self.id}, name={self.name}, contact_id={self.contact}, created={self.created}, updated={self.updated}'

class Module(_Base):
    name: Mapped[str] = mapped_column(index=True, unique=True)
    admin_id: Mapped[id] = mapped_column(ForeignKey('group.id',ondelete="CASCADE"))
    admin: Mapped[Group] = relationship(back_populates='managed_modules')
    questions: Mapped[list['Question']] = relationship(back_populates='module',  cascade='all, delete-orphan', passive_deletes=True)
    user_association: Mapped[list[AssociationModuleUser]] = relationship(back_populates='module')

    @classmethod
    def find(cls, name:str):
        return db.session.query(cls).filter_by(name=name)

    def __repr__(self):
        return f'id={self.id}, name={self.name}, admin_id={self.admin_id}, created={self.created}, updated={self.updated}'

class Question(_Base):
    module_id: Mapped[int] = mapped_column(ForeignKey('module.id',ondelete="CASCADE"))
    name: Mapped[str] = mapped_column()
    max_try: Mapped[int|None] = mapped_column()
    grader: Mapped[dict] = mapped_column(DillPickleType)
    module: Mapped[Module] = relationship(back_populates='questions')
    scores: Mapped[list['Score']] = relationship(back_populates='question', cascade='all, delete-orphan', passive_deletes=True)

    @classmethod
    def find(cls, module_id:int, name:str):
        return db.session.query(cls).filter_by(module_id=module_id,name=name)

    __table_args__ = (
        Index('ix_question_module_name', 'module_id', 'name', unique=True),
    )

    def __repr__(self):
        return f'id={self.id}, module={self.module_id}, name={self.name}, grader={self.grader}, created={self.created}, updated={self.updated}'

class Answer(_Base):
    sender_id: Mapped[int] = mapped_column(ForeignKey('user.id'))
    group_id: Mapped[int] = mapped_column(ForeignKey('group.id',ondelete="CASCADE"))
    score_id: Mapped[int] =  mapped_column(ForeignKey('score.id',ondelete="CASCADE"))
    score: Mapped['Score'] = relationship(back_populates='answers')
    group: Mapped['Group'] = relationship(back_populates='answers')

    __table_args__ = (
        Index('ix_answer_groupe_score', 'group_id', 'score_id', unique=True),
    )

    def __repr__(self):
        return f'id={self.id}, group_id={self.group_id}, sender_id={self.sender_id}, score_id={self.score_id}, created={self.created}, updated={self.updated}'

class Score(_Base):
    question_id: Mapped[int] = mapped_column(ForeignKey('question.id',ondelete="CASCADE"))
    digest: Mapped[str] = mapped_column(unique=True)
    data: Mapped[dict] = mapped_column(DillPickleType)
    score: Mapped[float|None] = mapped_column()
    question: Mapped['Question']  = relationship(back_populates='scores')
    answers: Mapped[list['Answer']]  = relationship(back_populates='score')

    def __repr__(self):
        return f'id={self.id}, question_id={self.question_id}, digest={self.digest}, data={self.data}, score={self.score}, created={self.created}, updated={self.updated}'

