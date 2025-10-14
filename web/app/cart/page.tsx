'use client'

import { useRouter } from 'next/navigation'
import { useCartStore, useAuthStore } from '@/lib/store'

export default function CartPage() {
  const router = useRouter()
  const { user } = useAuthStore()
  const { items, removeItem, updateQuantity, getTotalPrice, clearCart } = useCartStore()

  const formatPrice = (amount: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(amount)
  }

  const handleCheckout = () => {
    if (!user) {
      router.push('/login')
      return
    }
    // 決済ページに遷移（今後実装）
    router.push('/checkout')
  }

  if (items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="text-center">
          <h1 className="text-3xl font-bold mb-4">カート</h1>
          <p className="text-gray-600 mb-8">カートは空です</p>
          <button
            onClick={() => router.push('/products')}
            className="bg-primary text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
          >
            商品を見る
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">カート</h1>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* カート商品一覧 */}
        <div className="lg:col-span-2 space-y-4">
          {items.map((item) => (
            <div
              key={item.product.id}
              className="flex gap-4 border rounded-lg p-4"
            >
              <div className="w-24 h-24 bg-gray-200 rounded flex-shrink-0 flex items-center justify-center">
                <span className="text-xs text-gray-400">画像</span>
              </div>

              <div className="flex-1">
                <h3 className="font-semibold mb-1">{item.product.name}</h3>
                <p className="text-sm text-gray-600 mb-2">
                  {formatPrice(item.product.price.amount)}
                </p>

                <div className="flex items-center gap-2">
                  <button
                    onClick={() =>
                      updateQuantity(item.product.id, Math.max(1, item.quantity - 1))
                    }
                    className="w-8 h-8 border rounded hover:bg-gray-100"
                  >
                    -
                  </button>
                  <span className="w-12 text-center">{item.quantity}</span>
                  <button
                    onClick={() =>
                      updateQuantity(
                        item.product.id,
                        Math.min(item.product.stock_quantity, item.quantity + 1)
                      )
                    }
                    className="w-8 h-8 border rounded hover:bg-gray-100"
                  >
                    +
                  </button>
                </div>
              </div>

              <div className="text-right">
                <div className="font-bold mb-2">
                  {formatPrice(item.product.price.amount * item.quantity)}
                </div>
                <button
                  onClick={() => removeItem(item.product.id)}
                  className="text-sm text-red-600 hover:underline"
                >
                  削除
                </button>
              </div>
            </div>
          ))}

          <button
            onClick={clearCart}
            className="text-sm text-gray-600 hover:underline"
          >
            カートを空にする
          </button>
        </div>

        {/* 合計金額・購入ボタン */}
        <div className="lg:col-span-1">
          <div className="border rounded-lg p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-4">ご注文内容</h2>

            <div className="space-y-2 mb-4">
              <div className="flex justify-between">
                <span>小計</span>
                <span>{formatPrice(getTotalPrice())}</span>
              </div>
              <div className="flex justify-between">
                <span>配送料</span>
                <span>無料</span>
              </div>
            </div>

            <div className="border-t pt-4 mb-6">
              <div className="flex justify-between text-xl font-bold">
                <span>合計</span>
                <span className="text-primary">{formatPrice(getTotalPrice())}</span>
              </div>
            </div>

            <button
              onClick={handleCheckout}
              className="w-full bg-primary text-white py-3 rounded-lg font-semibold hover:bg-blue-600 transition"
            >
              購入手続きへ
            </button>

            <button
              onClick={() => router.push('/products')}
              className="w-full mt-2 border border-gray-300 py-3 rounded-lg font-semibold hover:bg-gray-50 transition"
            >
              買い物を続ける
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
