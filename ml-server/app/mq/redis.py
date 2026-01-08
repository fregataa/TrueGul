import asyncio
import logging

import redis.asyncio as redis

from app.mq.base import Consumer, Message, MessageHandler

logger = logging.getLogger(__name__)


class RedisConsumer(Consumer):
    def __init__(
        self,
        redis_url: str,
        stream_name: str,
        consumer_group: str,
        consumer_name: str,
    ):
        self._redis_url = redis_url
        self._stream_name = stream_name
        self._consumer_group = consumer_group
        self._consumer_name = consumer_name
        self._redis: redis.Redis | None = None
        self._running = False

    async def connect(self) -> None:
        client = redis.from_url(self._redis_url, decode_responses=True)
        self._redis = client
        try:
            await client.xgroup_create(
                self._stream_name, self._consumer_group, id="0", mkstream=True
            )
            logger.info(f"Created consumer group: {self._consumer_group}")
        except redis.ResponseError as e:
            if "BUSYGROUP" in str(e):
                logger.info(f"Consumer group {self._consumer_group} already exists")
            else:
                raise

    async def disconnect(self) -> None:
        if self._redis:
            await self._redis.aclose()

    async def start(self, handler: MessageHandler) -> None:
        if self._redis is None:
            raise RuntimeError("Redis client not connected. Call connect() first.")

        self._running = True
        logger.info("Starting Redis consumer...")
        client = self._redis

        while self._running:
            try:
                messages = await client.xreadgroup(
                    self._consumer_group,
                    self._consumer_name,
                    {self._stream_name: ">"},
                    count=1,
                    block=5000,
                )

                if not messages:
                    continue

                for _stream, stream_messages in messages:
                    for message_id, data in stream_messages:
                        message = Message(id=message_id, data=data)
                        await handler(message)

            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.exception(f"Error in consumer loop: {e}")
                await asyncio.sleep(1)

    async def stop(self) -> None:
        self._running = False

    async def ack(self, message_id: str) -> None:
        if self._redis:
            await self._redis.xack(self._stream_name, self._consumer_group, message_id)
