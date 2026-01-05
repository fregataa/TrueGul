"use client";

import { XCircle } from "lucide-react";
import Link from "next/link";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

interface AnalysisErrorProps {
  errorCode?: string;
  errorMessage?: string;
}

const ERROR_MESSAGES: Record<string, string> = {
  ML_MODEL_ERROR: "AI 분석 중 오류가 발생했습니다. 잠시 후 다시 시도해주세요.",
  OPENAI_API_ERROR: "피드백 생성 중 오류가 발생했습니다.",
  INVALID_INPUT: "글 내용이 분석하기에 적합하지 않습니다.",
  TIMEOUT: "분석 시간이 초과되었습니다.",
  INTERNAL_ERROR: "서버 오류가 발생했습니다.",
};

export function AnalysisError({ errorCode, errorMessage }: AnalysisErrorProps) {
  const message = errorCode
    ? ERROR_MESSAGES[errorCode] || errorMessage || "분석 중 오류가 발생했습니다."
    : errorMessage || "분석 중 오류가 발생했습니다.";

  return (
    <Alert variant="destructive">
      <XCircle className="h-4 w-4" />
      <AlertTitle>분석 실패</AlertTitle>
      <AlertDescription className="mt-2">
        <p>{message}</p>
        <div className="mt-4">
          <Button asChild variant="outline" size="sm">
            <Link href="/dashboard">목록으로 돌아가기</Link>
          </Button>
        </div>
      </AlertDescription>
    </Alert>
  );
}
