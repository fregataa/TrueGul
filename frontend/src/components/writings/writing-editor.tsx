"use client";

import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import type { Writing, WritingType } from "@/lib/api/writings";

const MAX_CONTENT_LENGTH = 2000;

interface WritingEditorProps {
  mode: "create" | "edit";
  initialData?: Writing;
  onSave: (type: WritingType, title: string, content: string) => Promise<Writing | null>;
  isLoading: boolean;
  error: string | null;
}

export function WritingEditor({ mode, initialData, onSave, isLoading, error }: WritingEditorProps) {
  const router = useRouter();
  const [type, setType] = useState<WritingType>(initialData?.type || "essay");
  const [title, setTitle] = useState(initialData?.title || "");
  const [content, setContent] = useState(initialData?.content || "");
  const [validationError, setValidationError] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);

  const contentLength = [...content].length;
  const isOverLimit = contentLength > MAX_CONTENT_LENGTH;

  const AUTO_SAVE_INTERVAL_MS = 30000;

  const autoSave = useCallback(async () => {
    if (mode === "edit" && initialData && title && content && !isOverLimit) {
      const hasChanges =
        title !== initialData.title || content !== initialData.content || type !== initialData.type;

      if (hasChanges && !isSaving) {
        setIsSaving(true);
        const result = await onSave(type, title, content);
        setIsSaving(false);
        if (result) {
          setLastSaved(new Date());
        }
      }
    }
  }, [mode, initialData, title, content, type, isOverLimit, isSaving, onSave]);

  useEffect(() => {
    if (mode === "edit") {
      const interval = setInterval(autoSave, AUTO_SAVE_INTERVAL_MS);
      return () => clearInterval(interval);
    }
  }, [mode, autoSave]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setValidationError(null);

    if (!title.trim()) {
      setValidationError("제목을 입력해주세요.");
      return;
    }

    if (!content.trim()) {
      setValidationError("내용을 입력해주세요.");
      return;
    }

    if (isOverLimit) {
      setValidationError(`내용은 ${MAX_CONTENT_LENGTH}자를 초과할 수 없습니다.`);
      return;
    }

    const result = await onSave(type, title, content);
    if (result) {
      router.push("/dashboard");
    }
  };

  const handleManualSave = async () => {
    if (!title.trim() || !content.trim() || isOverLimit) return;

    setIsSaving(true);
    const result = await onSave(type, title, content);
    setIsSaving(false);

    if (result) {
      setLastSaved(new Date());
    }
  };

  const displayError = validationError || error;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {displayError && (
        <div className="p-4 text-sm text-red-500 bg-red-50 rounded-md">{displayError}</div>
      )}

      <fieldset>
        <legend className="block text-sm font-medium mb-2">글 종류</legend>
        <div className="flex gap-4">
          <label htmlFor="type-essay" className="flex items-center gap-2 cursor-pointer">
            <input
              id="type-essay"
              type="radio"
              name="type"
              value="essay"
              checked={type === "essay"}
              onChange={(e) => setType(e.target.value as WritingType)}
              className="w-4 h-4 text-blue-600"
            />
            <span>에세이</span>
          </label>
          <label htmlFor="type-cover-letter" className="flex items-center gap-2 cursor-pointer">
            <input
              id="type-cover-letter"
              type="radio"
              name="type"
              value="cover_letter"
              checked={type === "cover_letter"}
              onChange={(e) => setType(e.target.value as WritingType)}
              className="w-4 h-4 text-blue-600"
            />
            <span>자기소개서</span>
          </label>
        </div>
      </fieldset>

      <div>
        <label htmlFor="title" className="block text-sm font-medium mb-2">
          제목
        </label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          maxLength={255}
          className="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="글 제목을 입력하세요"
        />
      </div>

      <div>
        <div className="flex justify-between items-center mb-2">
          <label htmlFor="content" className="block text-sm font-medium">
            내용
          </label>
          <span
            className={`text-sm ${isOverLimit ? "text-red-500 font-semibold" : "text-gray-500"}`}
          >
            {contentLength} / {MAX_CONTENT_LENGTH}자
          </span>
        </div>
        <textarea
          id="content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={15}
          className={`w-full px-4 py-3 border rounded-md focus:outline-none focus:ring-2 resize-none ${
            isOverLimit ? "border-red-500 focus:ring-red-500" : "focus:ring-blue-500"
          }`}
          placeholder="글 내용을 입력하세요"
        />
        {isOverLimit && (
          <p className="mt-1 text-sm text-red-500">
            글자 수가 {MAX_CONTENT_LENGTH}자를 초과했습니다. ({contentLength - MAX_CONTENT_LENGTH}자
            초과)
          </p>
        )}
      </div>

      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          <Button type="submit" disabled={isLoading || isOverLimit}>
            {isLoading ? "저장 중..." : mode === "create" ? "작성 완료" : "수정 완료"}
          </Button>
          {mode === "edit" && (
            <Button
              type="button"
              variant="outline"
              onClick={handleManualSave}
              disabled={isSaving || isOverLimit || !title.trim() || !content.trim()}
            >
              {isSaving ? "저장 중..." : "임시 저장"}
            </Button>
          )}
          <Button type="button" variant="outline" onClick={() => router.push("/dashboard")}>
            취소
          </Button>
        </div>

        {lastSaved && (
          <span className="text-sm text-gray-500">
            마지막 저장: {lastSaved.toLocaleTimeString("ko-KR")}
          </span>
        )}
      </div>
    </form>
  );
}
