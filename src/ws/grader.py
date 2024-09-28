from abc import ABC, abstractmethod

class Grader(ABC):

    @abstractmethod
    def validate(self):
        pass

    @abstractmethod
    def score(self, answer : bytes) -> dict:
        pass

class Factory(ABC):
    @abstractmethod
    def new(self) -> Grader:
        pass
