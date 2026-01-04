"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect } from "react";
import { WritingEditor } from "@/components/writings/writing-editor";
import type { WritingType } from "@/lib/api/writings";
import { useAuthStore } from "@/stores/auth";
import { useWritingsStore } from "@/stores/writings";

export default function EditWritingPage() {
  const router = useRouter();
  const params = useParams();
  const id = params.id as string;

  const { isAuthenticated, isLoading: authLoading, checkAuth } = useAuthStore();
  const { currentWriting, isLoading, error, fetchWriting, updateWriting } = useWritingsStore();

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
    }
  }, [isAuthenticated, id, fetchWriting]);

  const handleSave = async (type: WritingType, title: string, content: string) => {
    return updateWriting(id, { type, title, content });
  };

  if (authLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (isLoading && !currentWriting) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>Loading...</p>
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
        <div className="max-w-4xl mx-auto px-4 py-4">
          <h1 className="text-xl font-bold">글 수정</h1>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <WritingEditor
            mode="edit"
            initialData={currentWriting}
            onSave={handleSave}
            isLoading={isLoading}
            error={error}
          />
        </div>
      </main>
    </div>
  );
}
