"use client";

import { AlertTriangle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

interface RateLimitWarningProps {
  message?: string;
}

export function RateLimitWarning({ message }: RateLimitWarningProps) {
  return (
    <Alert variant="warning">
      <AlertTriangle className="h-4 w-4" />
      <AlertTitle>제출 횟수 초과</AlertTitle>
      <AlertDescription>
        {message || "오늘의 제출 횟수를 모두 사용했습니다. 내일 다시 시도해주세요."}
      </AlertDescription>
    </Alert>
  );
}
