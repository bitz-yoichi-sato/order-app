import Link from "next/link";

export default function Home() {
  return (
    <main className="page">
      <div className="page__eyebrow">Order Management</div>
      <h1 className="page__title">受注管理システム</h1>
      <p className="page__lead">
        受注データの一覧・詳細を確認できる社内向け管理画面です。
      </p>
      <Link href="/orders" className="button">
        注文一覧を見る →
      </Link>
    </main>
  );
}
