"use client";

import { CheckCircle, Clock } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Analysis } from "@/lib/api/analysis";
import { AIScoreBadge } from "./ai-score-badge";

interface AnalysisResultProps {
  analysis: Analysis;
}

export function AnalysisResult({ analysis }: AnalysisResultProps) {
  const aiScore = analysis.ai_score ?? 0;
  const latencySeconds = analysis.latency_ms ? (analysis.latency_ms / 1000).toFixed(1) : null;

  return (
    <Card>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <CheckCircle className="h-5 w-5 text-green-600" />
            <span>분석 결과</span>
          </div>
          <AIScoreBadge score={aiScore} size="lg" />
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {analysis.feedback && (
          <div>
            <h4 className="text-sm font-medium text-gray-700 mb-2">피드백</h4>
            <div className="rounded-lg bg-gray-50 p-4 text-sm text-gray-700 whitespace-pre-wrap leading-relaxed">
              {analysis.feedback}
            </div>
          </div>
        )}

        {latencySeconds && (
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <Clock className="h-3.5 w-3.5" />
            <span>분석 소요 시간: {latencySeconds}초</span>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
