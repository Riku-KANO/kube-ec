'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { productApi } from '@/lib/api'
import type { Product } from '@/lib/types'

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')

  useEffect(() => {
    loadProducts()
  }, [category, search])

  const loadProducts = async () => {
    try {
      setLoading(true)
      const data = await productApi.list(1, 20, category, search)
      setProducts(data.products || [])
    } catch (err) {
      setError('商品の読み込みに失敗しました')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const formatPrice = (amount: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(amount)
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">読み込み中...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-red-500">{error}</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">商品一覧</h1>

      {/* 検索・フィルター */}
      <div className="mb-8 flex gap-4">
        <input
          type="text"
          placeholder="商品を検索..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 px-4 py-2 border rounded-lg"
        />
        <select
          value={category}
          onChange={(e) => setCategory(e.target.value)}
          className="px-4 py-2 border rounded-lg"
        >
          <option value="">すべてのカテゴリー</option>
          <option value="electronics">家電</option>
          <option value="fashion">ファッション</option>
          <option value="food">食品</option>
          <option value="books">書籍</option>
        </select>
      </div>

      {/* 商品グリッド */}
      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((product) => (
          <Link
            key={product.id}
            href={`/products/${product.id}`}
            className="border rounded-lg p-4 hover:shadow-lg transition"
          >
            <div className="aspect-square bg-gray-200 rounded-lg mb-4 flex items-center justify-center">
              <span className="text-gray-400">画像</span>
            </div>
            <h3 className="font-semibold mb-2 truncate">{product.name}</h3>
            <p className="text-gray-600 text-sm mb-2 line-clamp-2">
              {product.description}
            </p>
            <div className="flex justify-between items-center">
              <span className="text-lg font-bold text-primary">
                {formatPrice(product.price.amount)}
              </span>
              <span className="text-sm text-gray-500">
                在庫: {product.stock_quantity}
              </span>
            </div>
          </Link>
        ))}
      </div>

      {products.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          商品が見つかりませんでした
        </div>
      )}
    </div>
  )
}
