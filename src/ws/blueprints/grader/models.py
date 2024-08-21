import datetime
from ws import db, DillPickleType
from sqlalchemy import Integer, String, ForeignKey, LargeBinary, Index, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column

class _Base(db.Model):
    __abstract__ = True
    id : Mapped[int] = mapped_column(primary_key=True)
    created: Mapped[datetime.datetime] = mapped_column(default=datetime.datetime.utcnow, nullable=False)
    updated: Mapped[datetime.datetime] = mapped_column(default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow, nullable=False)


class User(_Base):
    username: Mapped[str] = mapped_column(index=True, unique=True, nullable=False)
    email: Mapped[str] = mapped_column(unique=True, nullable=False)
    public_key: Mapped[bytes] = mapped_column(LargeBinary, unique=True, nullable=False)
    answers: Mapped[list['Answer']] = db.relationship('Answer', back_populates='user', cascade='all, delete-orphan', passive_deletes=True)

    def __repr__(self):
        return f'id={self.id}, username={self.username}, key={self.public_key}, created={self.created}, updated={self.updated}'


class Question(_Base):
    module: Mapped[str] = mapped_column(index=True, nullable=False)
    name: Mapped[str] = mapped_column(nullable=False)
    grader: Mapped[dict] = mapped_column(DillPickleType, nullable=False)
    answers: Mapped[list['Answer']] = db.relationship('Answer', back_populates='question', cascade='all, delete-orphan', passive_deletes=True)

    __table_args__ = (
        Index('ix_module_name', 'module', 'name', unique=True),
    )

    def __repr__(self):
        return f'id={self.id}, module={self.module}, name={self.name}, grader={self.grader}, created={self.created}, updated={self.updated}'


class Answer(_Base):
    user_id: Mapped[int] = mapped_column(ForeignKey('user.id',ondelete="CASCADE"), nullable=False)
    question_id: Mapped[int] = mapped_column(ForeignKey('question.id',ondelete="CASCADE"), nullable=False)
    score_id: Mapped[int] =  mapped_column(ForeignKey('score.id',ondelete="CASCADE"), nullable=False)
    group_name: Mapped[str] = mapped_column(nullable=True)
    user: Mapped['User'] = db.relationship('User', back_populates='answers')
    question: Mapped['Question']  = db.relationship('Question', back_populates='answers')
    score: Mapped['Score'] = db.relationship('Score', back_populates='answers')

    def __repr__(self):
        return f'id={self.id}, user_id={self.user_id}, question_id={self.question_id}, score_id={self.score_id}, created={self.created}, updated={self.updated}'

class Score(_Base):
    digest: Mapped[str] = mapped_column(unique=True, nullable=False)
    data: Mapped[dict] = mapped_column(DillPickleType, nullable=False)
    score: Mapped[float] = mapped_column(unique=False, nullable=True)
    answers: Mapped[list['Answer']] = db.relationship('Answer', back_populates='score', cascade='all, delete-orphan', passive_deletes=True)

    def __repr__(self):
        return f'id={self.id}, digest={self.digest}, data={self.data}, score={self.score}, created={self.created}, updated={self.updated}'

