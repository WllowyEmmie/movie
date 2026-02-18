import type { Metadata } from "next";
import { Montserrat, Roboto } from "next/font/google";
import "./globals.css";

/* Load Montserrat as main font */
const montserrat = Montserrat({
  variable: "--font-sans",
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700", "800", "900"], // add all weights if needed
});
const roboto = Roboto({
  variable: "--font-sans",
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700", "800", "900"], // add all weights if needed
});

export const metadata: Metadata = {
  title: "Movie Recommender",
  description: "AI-powered movie recommendation app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={montserrat.variable}>
      <body className="antialiased bg-background text-foreground">
        {children}
      </body>
    </html>
  );
}