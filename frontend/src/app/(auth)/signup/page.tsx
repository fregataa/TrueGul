"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { AuthForm } from "@/components/auth/auth-form";
import { useAuthStore } from "@/stores/auth";

export default function SignupPage() {
  const router = useRouter();
  const { signup, isLoading, isAuthenticated, error, clearError } = useAuthStore();

  useEffect(() => {
    clearError();
  }, [clearError]);

  useEffect(() => {
    if (isAuthenticated) {
      router.push("/dashboard");
    }
  }, [isAuthenticated, router]);

  const handleSubmit = async (email: string, password: string) => {
    const success = await signup(email, password);
    if (success) {
      router.push("/login?registered=true");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-8 bg-white rounded-lg shadow-md">
        <h1 className="text-2xl font-bold text-center mb-6">Sign Up</h1>

        <AuthForm mode="signup" onSubmit={handleSubmit} isLoading={isLoading} error={error} />

        <p className="mt-4 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <Link href="/login" className="text-blue-600 hover:underline">
            Log in
          </Link>
        </p>
      </div>
    </div>
  );
}
