"use client";

import { Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { Writing } from "@/lib/api/writings";
import { useAuthStore } from "@/stores/auth";
import { useWritingsStore } from "@/stores/writings";

export default function DashboardPage() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading: authLoading, logout, checkAuth } = useAuthStore();
  const {
    writings,
    total,
    page,
    totalPages,
    isLoading: writingsLoading,
    error,
    fetchWritings,
    deleteWriting,
  } = useWritingsStore();

  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  useEffect(() => {
    if (isAuthenticated) {
      fetchWritings();
    }
  }, [isAuthenticated, fetchWritings]);

  const handleLogout = async () => {
    await logout();
    router.push("/login");
  };

  const handleDelete = async (id: string) => {
    const success = await deleteWriting(id);
    if (success) {
      setDeleteConfirm(null);
    }
  };

  const handlePageChange = (newPage: number) => {
    fetchWritings(newPage);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("ko-KR", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getTypeLabel = (type: string) => {
    return type === "essay" ? "에세이" : "자소서";
  };

  const getStatusConfig = (status: string) => {
    switch (status) {
      case "draft":
        return { label: "임시저장", variant: "warning" as const };
      case "submitted":
        return { label: "분석 중", variant: "secondary" as const, showLoader: true };
      case "analyzed":
        return { label: "분석완료", variant: "success" as const };
      default:
        return { label: status, variant: "outline" as const };
    }
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
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">TrueGul</h1>
          <div className="flex items-center gap-4">
            <span className="text-sm text-gray-600">{user?.email}</span>
            <Button variant="outline" onClick={handleLogout}>
              Log Out
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">내 글 목록</h2>
          <Link href="/writings/new">
            <Button>새 글 작성</Button>
          </Link>
        </div>

        {error && <div className="p-4 mb-4 text-sm text-red-500 bg-red-50 rounded-md">{error}</div>}

        {writingsLoading ? (
          <div className="text-center py-8">
            <p className="text-gray-500">로딩 중...</p>
          </div>
        ) : writings.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <p className="text-gray-500 mb-4">작성한 글이 없습니다.</p>
            <Link href="/writings/new">
              <Button>첫 글 작성하기</Button>
            </Link>
          </div>
        ) : (
          <>
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      제목
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      종류
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      상태
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      수정일
                    </th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      작업
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {writings.map((writing: Writing) => (
                    <tr key={writing.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Link
                          href={`/writings/${writing.id}`}
                          className="text-blue-600 hover:underline font-medium"
                        >
                          {writing.title}
                        </Link>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {getTypeLabel(writing.type)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {(() => {
                          const config = getStatusConfig(writing.status);
                          return (
                            <Badge variant={config.variant} className="gap-1">
                              {config.showLoader && <Loader2 className="h-3 w-3 animate-spin" />}
                              {config.label}
                            </Badge>
                          );
                        })()}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {formatDate(writing.updated_at)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <div className="flex justify-end gap-2">
                          {writing.status === "draft" && (
                            <Link href={`/writings/${writing.id}/edit`}>
                              <Button variant="outline" size="sm">
                                수정
                              </Button>
                            </Link>
                          )}
                          {deleteConfirm === writing.id ? (
                            <div className="flex gap-1">
                              <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => handleDelete(writing.id)}
                              >
                                확인
                              </Button>
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => setDeleteConfirm(null)}
                              >
                                취소
                              </Button>
                            </div>
                          ) : (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => setDeleteConfirm(writing.id)}
                            >
                              삭제
                            </Button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {totalPages > 1 && (
              <div className="mt-4 flex justify-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => handlePageChange(page - 1)}
                >
                  이전
                </Button>
                <span className="flex items-center px-4 text-sm text-gray-600">
                  {page} / {totalPages} 페이지 (총 {total}개)
                </span>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages}
                  onClick={() => handlePageChange(page + 1)}
                >
                  다음
                </Button>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
