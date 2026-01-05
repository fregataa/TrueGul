"use client";

import { AlertTriangle, Send } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface SubmitConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void;
  isLoading: boolean;
  writingTitle: string;
}

export function SubmitConfirmDialog({
  open,
  onOpenChange,
  onConfirm,
  isLoading,
  writingTitle,
}: SubmitConfirmDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Send className="h-5 w-5" />
            분석 제출
          </DialogTitle>
          <DialogDescription className="text-left">
            <span className="font-medium text-gray-900">"{writingTitle}"</span>을(를) AI 분석에
            제출하시겠습니까?
          </DialogDescription>
        </DialogHeader>

        <div className="rounded-lg bg-yellow-50 p-4 border border-yellow-200">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-5 w-5 text-yellow-600 mt-0.5 flex-shrink-0" />
            <div className="text-sm text-yellow-800">
              <p className="font-medium mb-1">주의사항</p>
              <ul className="list-disc list-inside space-y-1 text-yellow-700">
                <li>하루 최대 5회까지 분석을 제출할 수 있습니다.</li>
                <li>제출 후에는 글을 수정할 수 없습니다.</li>
                <li>분석은 보통 10초 ~ 1분 정도 소요됩니다.</li>
              </ul>
            </div>
          </div>
        </div>

        <DialogFooter className="gap-2 sm:gap-0">
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
            취소
          </Button>
          <Button onClick={onConfirm} disabled={isLoading}>
            {isLoading ? "제출 중..." : "제출하기"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
