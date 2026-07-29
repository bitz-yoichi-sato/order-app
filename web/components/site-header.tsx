import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="site-header">
      <div className="site-header__inner">
        <Link href="/" className="site-header__brand">
          受注管理システム
        </Link>
        <nav className="site-header__nav">
          <Link href="/orders">注文一覧</Link>
        </nav>
      </div>
    </header>
  );
}
