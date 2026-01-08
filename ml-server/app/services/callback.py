import logging

import httpx

from app.config import settings
from app.schemas.task import AnalysisCallback

logger = logging.getLogger(__name__)

CALLBACK_HEADER = "X-Callback-Secret"


class CallbackClient:
    def __init__(self):
        self._client = httpx.AsyncClient(timeout=30.0)

    async def send_callback(self, callback_url: str, payload: AnalysisCallback) -> bool:
        headers = {CALLBACK_HEADER: settings.ml_callback_secret}

        try:
            response = await self._client.post(
                callback_url,
                json=payload.model_dump(),
                headers=headers,
            )
            response.raise_for_status()
            logger.info(f"Callback sent successfully for task {payload.task_id}")
            return True
        except httpx.HTTPError as e:
            logger.error(f"Callback failed for task {payload.task_id}: {e}")
            return False

    async def close(self):
        await self._client.aclose()
