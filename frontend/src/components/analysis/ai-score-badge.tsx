"use client";

import { cn } from "@/lib/utils";

interface AIScoreBadgeProps {
  score: number;
  size?: "sm" | "md" | "lg";
}

function getScoreConfig(score: number) {
  if (score <= 30) {
    return {
      label: "안전",
      bgColor: "bg-green-100",
      textColor: "text-green-800",
      borderColor: "border-green-200",
    };
  }
  if (score <= 60) {
    return {
      label: "주의",
      bgColor: "bg-yellow-100",
      textColor: "text-yellow-800",
      borderColor: "border-yellow-200",
    };
  }
  return {
    label: "위험",
    bgColor: "bg-red-100",
    textColor: "text-red-800",
    borderColor: "border-red-200",
  };
}

const sizeStyles = {
  sm: "px-2 py-0.5 text-xs",
  md: "px-3 py-1 text-sm",
  lg: "px-4 py-1.5 text-base",
};

export function AIScoreBadge({ score, size = "md" }: AIScoreBadgeProps) {
  const config = getScoreConfig(score);
  const roundedScore = Math.round(score);

  return (
    <div
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border font-semibold",
        config.bgColor,
        config.textColor,
        config.borderColor,
        sizeStyles[size]
      )}
    >
      <span>AI 확률</span>
      <span className="font-bold">{roundedScore}%</span>
      <span className="text-xs font-normal">({config.label})</span>
    </div>
  );
}
