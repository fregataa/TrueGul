"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { WritingEditor } from "@/components/writings/writing-editor";
import type { WritingType } from "@/lib/api/writings";
import { useAuthStore } from "@/stores/auth";
import { useWritingsStore } from "@/stores/writings";

export default function NewWritingPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading, checkAuth } = useAuthStore();
  const { createWriting, isLoading, error } = useWritingsStore();

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  const handleSave = async (type: WritingType, title: string, content: string) => {
    return createWriting(type, title, content);
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

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <h1 className="text-xl font-bold">새 글 작성</h1>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <WritingEditor mode="create" onSave={handleSave} isLoading={isLoading} error={error} />
        </div>
      </main>
    </div>
  );
}
