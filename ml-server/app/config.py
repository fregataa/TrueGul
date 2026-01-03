from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    openai_api_key: str = ""
    model_version: str = "0.1.0"

    class Config:
        env_file = ".env"


settings = Settings()
