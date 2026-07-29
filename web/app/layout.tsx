import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "受注管理システム",
  description: "受注管理システム 社内向け画面",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
