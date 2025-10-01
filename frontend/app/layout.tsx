import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Navigation from "../components/Navigation";
import { AuthProvider } from "@/contexts/AuthContext";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Recipe API - Find and Save Your Favorite Recipes",
  description: "Discover recipes, find what to cook with ingredients you have, and save your favorites",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased h-full bg-gray-50`}>
        <AuthProvider>
          <div className="flex flex-col h-full">
            <Navigation />
            <main className="flex-1 bg-gray-50">
              {children}
            </main>
          </div>
        </AuthProvider>
      </body>
    </html>
  );

}