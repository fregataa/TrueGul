from abc import ABC, abstractmethod
from dataclasses import dataclass
from typing import Callable, Awaitable


@dataclass
class Message:
    id: str
    data: dict


MessageHandler = Callable[[Message], Awaitable[None]]


class Consumer(ABC):
    @abstractmethod
    async def connect(self) -> None:
        pass

    @abstractmethod
    async def disconnect(self) -> None:
        pass

    @abstractmethod
    async def start(self, handler: MessageHandler) -> None:
        pass

    @abstractmethod
    async def stop(self) -> None:
        pass

    @abstractmethod
    async def ack(self, message_id: str) -> None:
        pass
