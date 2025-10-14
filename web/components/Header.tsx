'use client'

import Link from 'next/link'
import { useAuthStore, useCartStore } from '@/lib/store'

export default function Header() {
  const { user, logout } = useAuthStore()
  const { getTotalItems } = useCartStore()

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <Link href="/" className="text-2xl font-bold text-primary">
            Kube EC
          </Link>

          <nav className="flex items-center space-x-6">
            <Link href="/products" className="hover:text-primary transition">
              商品一覧
            </Link>

            {user ? (
              <>
                <Link href="/orders" className="hover:text-primary transition">
                  注文履歴
                </Link>
                <Link href="/cart" className="hover:text-primary transition relative">
                  カート
                  {getTotalItems() > 0 && (
                    <span className="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                      {getTotalItems()}
                    </span>
                  )}
                </Link>
                <button
                  onClick={logout}
                  className="hover:text-primary transition"
                >
                  ログアウト
                </button>
                <span className="text-gray-600">
                  {user.name} さん
                </span>
              </>
            ) : (
              <>
                <Link href="/login" className="hover:text-primary transition">
                  ログイン
                </Link>
                <Link
                  href="/register"
                  className="bg-primary text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                >
                  新規登録
                </Link>
              </>
            )}
          </nav>
        </div>
      </div>
    </header>
  )
}
