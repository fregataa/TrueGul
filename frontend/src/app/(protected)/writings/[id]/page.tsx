"use client";

import { Send } from "lucide-react";
import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { AnalysisError } from "@/components/analysis/analysis-error";
import { AnalysisProgress } from "@/components/analysis/analysis-progress";
import { AnalysisResult } from "@/components/analysis/analysis-result";
import { RateLimitWarning } from "@/components/analysis/rate-limit-warning";
import { SubmitConfirmDialog } from "@/components/analysis/submit-confirm-dialog";
import { Button } from "@/components/ui/button";
import { useAnalysisStore } from "@/stores/analysis";
import { useAuthStore } from "@/stores/auth";
import { useWritingsStore } from "@/stores/writings";

export default function ViewWritingPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;

  const [showSubmitDialog, setShowSubmitDialog] = useState(false);

  const { isAuthenticated, isLoading: authLoading, checkAuth } = useAuthStore();
  const { currentWriting, isLoading, error, fetchWriting } = useWritingsStore();
  const {
    currentAnalysis,
    isSubmitting,
    isPolling,
    isRateLimited,
    error: analysisError,
    submitWriting,
    fetchAnalysis,
    stopPolling,
    clearAnalysis,
  } = useAnalysisStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  useEffect(() => {
    if (isAuthenticated && id) {
      fetchWriting(id);
      fetchAnalysis(id);
    }
  }, [isAuthenticated, id, fetchWriting, fetchAnalysis]);

  useEffect(() => {
    return () => {
      stopPolling();
      clearAnalysis();
    };
  }, [stopPolling, clearAnalysis]);

  const handleSubmit = async () => {
    const success = await submitWriting(id);
    if (success) {
      setShowSubmitDialog(false);
      fetchWriting(id);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("ko-KR", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getTypeLabel = (type: string) => {
    return type === "essay" ? "에세이" : "자기소개서";
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case "draft":
        return "임시저장";
      case "submitted":
        return "제출됨";
      case "analyzed":
        return "분석완료";
      default:
        return status;
    }
  };

  const shouldShowSubmitButton =
    currentWriting?.status === "draft" && !isPolling && !currentAnalysis;

  const shouldShowProgress =
    isPolling ||
    (currentAnalysis &&
      (currentAnalysis.status === "pending" || currentAnalysis.status === "processing"));

  const shouldShowResult = currentAnalysis?.status === "completed";

  const shouldShowError = currentAnalysis?.status === "failed";

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <main className="max-w-4xl mx-auto px-4 py-8">
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-red-500 mb-4">{error}</p>
            <Link href="/dashboard">
              <Button>목록으로 돌아가기</Button>
            </Link>
          </div>
        </main>
      </div>
    );
  }

  if (!currentWriting) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>글을 찾을 수 없습니다.</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-4xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold truncate">{currentWriting.title}</h1>
          <div className="flex gap-2">
            {currentWriting.status === "draft" && (
              <Link href={`/writings/${id}/edit`}>
                <Button variant="outline">수정</Button>
              </Link>
            )}
            <Link href="/dashboard">
              <Button variant="outline">목록</Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8 space-y-6">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex gap-4 mb-6 text-sm text-gray-500">
            <span className="inline-flex items-center px-2 py-1 bg-gray-100 rounded">
              {getTypeLabel(currentWriting.type)}
            </span>
            <span
              className={`inline-flex items-center px-2 py-1 rounded ${
                currentWriting.status === "draft"
                  ? "bg-yellow-100 text-yellow-800"
                  : currentWriting.status === "submitted"
                    ? "bg-blue-100 text-blue-800"
                    : "bg-green-100 text-green-800"
              }`}
            >
              {getStatusLabel(currentWriting.status)}
            </span>
            <span>작성일: {formatDate(currentWriting.created_at)}</span>
            <span>수정일: {formatDate(currentWriting.updated_at)}</span>
          </div>

          <div className="prose max-w-none">
            <div className="whitespace-pre-wrap text-gray-800 leading-relaxed">
              {currentWriting.content}
            </div>
          </div>

          <div className="mt-6 pt-4 border-t text-sm text-gray-500">
            글자 수: {[...currentWriting.content].length}자
          </div>
        </div>

        {isRateLimited && <RateLimitWarning message={analysisError || undefined} />}

        {shouldShowSubmitButton && !isRateLimited && (
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex flex-col items-center text-center">
              <h3 className="text-lg font-semibold text-gray-900 mb-2">AI 분석 제출</h3>
              <p className="text-sm text-gray-500 mb-4">
                작성한 글을 AI에게 분석 요청하여 피드백을 받아보세요.
              </p>
              <Button onClick={() => setShowSubmitDialog(true)} className="gap-2">
                <Send className="h-4 w-4" />
                분석 제출하기
              </Button>
            </div>
          </div>
        )}

        {shouldShowProgress && currentAnalysis && (
          <div className="bg-white rounded-lg shadow">
            <AnalysisProgress status={currentAnalysis.status} />
          </div>
        )}

        {shouldShowResult && currentAnalysis && <AnalysisResult analysis={currentAnalysis} />}

        {shouldShowError && currentAnalysis && (
          <AnalysisError
            errorCode={currentAnalysis.error_code}
            errorMessage={currentAnalysis.error_message}
          />
        )}
      </main>

      <SubmitConfirmDialog
        open={showSubmitDialog}
        onOpenChange={setShowSubmitDialog}
        onConfirm={handleSubmit}
        isLoading={isSubmitting}
        writingTitle={currentWriting.title}
      />
    </div>
  );
}
