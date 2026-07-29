import Link from "next/link";

export default function Home() {
  return (
    <main style={{ padding: "2rem" }}>
      <h1>受注管理システム</h1>
      <p>
        <Link href="/orders">注文一覧へ</Link>
      </p>
    </main>
  );
}
