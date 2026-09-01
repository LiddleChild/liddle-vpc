import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Expenses",
  description: "View current monthly expenses across all accounts.",
  icons: {
    icon: "/icon.svg",
    apple: "/apple-icon.png",
  },
  appleWebApp: { capable: true, title: "Expenses", statusBarStyle: "black-translucent" },
};

export const viewport = {
  width: "device-width",
  initialScale: 1,
  viewportFit: "cover",
  themeColor: "#ffffff",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className="h-full antialiased">
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
