import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50">
      <main className="text-center px-4">
        <h1 className="text-4xl font-bold mb-4">TrueGul</h1>
        <p className="text-xl text-gray-600 mb-8">AI 없이 글쓰기 훈련을 돕는 서비스</p>

        <div className="flex gap-4 justify-center">
          <Link href="/login">
            <Button variant="outline">Log In</Button>
          </Link>
          <Link href="/signup">
            <Button>Sign Up</Button>
          </Link>
        </div>
      </main>
    </div>
  );
}
