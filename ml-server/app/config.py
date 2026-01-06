from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    port: int = 8000
    redis_url: str = "redis://localhost:6379"
    api_server_url: str = "http://localhost:8080"
    ml_callback_secret: str = ""
    model_name: str = "Hello-SimpleAI/chatgpt-detector-roberta"
    model_version: str = "0.1.0"
    feedback_model_path: str = "/app/models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf"
    feedback_n_ctx: int = 2048
    feedback_n_threads: int = 4
    feedback_max_tokens: int = 256
    stream_name: str = "analysis_tasks"
    consumer_group: str = "ml_workers"
    consumer_name: str = "worker_1"
    max_retries: int = 3

    class Config:
        env_file = ".env"


settings = Settings()
