"use client";

import { Clock, Cpu, Loader2 } from "lucide-react";
import type { AnalysisStatus } from "@/lib/api/analysis";

interface AnalysisProgressProps {
  status: AnalysisStatus;
}

const statusConfig: Record<
  AnalysisStatus,
  { icon: typeof Loader2; text: string; description: string }
> = {
  pending: {
    icon: Clock,
    text: "분석 대기 중",
    description: "곧 분석이 시작됩니다...",
  },
  processing: {
    icon: Cpu,
    text: "분석 중",
    description: "AI가 글을 분석하고 있습니다...",
  },
  completed: {
    icon: Loader2,
    text: "완료",
    description: "",
  },
  failed: {
    icon: Loader2,
    text: "실패",
    description: "",
  },
};

export function AnalysisProgress({ status }: AnalysisProgressProps) {
  const config = statusConfig[status];
  const Icon = config.icon;

  if (status === "completed" || status === "failed") {
    return null;
  }

  return (
    <div className="flex flex-col items-center justify-center py-12 px-4">
      <div className="relative">
        <div className="absolute inset-0 animate-ping rounded-full bg-blue-400 opacity-20" />
        <div className="relative flex h-16 w-16 items-center justify-center rounded-full bg-blue-100">
          {status === "processing" ? (
            <Loader2 className="h-8 w-8 text-blue-600 animate-spin" />
          ) : (
            <Icon className="h-8 w-8 text-blue-600" />
          )}
        </div>
      </div>
      <h3 className="mt-6 text-lg font-semibold text-gray-900">{config.text}</h3>
      <p className="mt-2 text-sm text-gray-500">{config.description}</p>
    </div>
  );
}
