'use client';

import { use, useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { productApi } from '@/lib/api';
import { useCartStore, useAuthStore } from '@/lib/store';
import type { Product } from '@/lib/types';

export default function ProductDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [product, setProduct] = useState<Product | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);
  const { addItem } = useCartStore();
  const { user } = useAuthStore();

  useEffect(() => {
    loadProduct();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const loadProduct = async () => {
    try {
      const data = await productApi.get(id);
      setProduct(data);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = () => {
    if (!product) return;
    addItem(product, quantity);
    alert('カートに追加しました');
  };

  const handleBuyNow = () => {
    if (!user) {
      router.push('/login');
      return;
    }
    if (!product) return;
    addItem(product, quantity);
    router.push('/cart');
  };

  const formatPrice = (amount: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(amount);
  };

  if (loading) {
    return <div className="container mx-auto px-4 py-8">読み込み中...</div>;
  }

  if (!product) {
    return <div className="container mx-auto px-4 py-8">商品が見つかりません</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid md:grid-cols-2 gap-8">
        {/* 商品画像 */}
        <div className="aspect-square bg-gray-200 rounded-lg flex items-center justify-center">
          <span className="text-gray-400 text-2xl">商品画像</span>
        </div>

        {/* 商品情報 */}
        <div>
          <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
          <div className="text-3xl font-bold text-primary mb-6">
            {formatPrice(product.price.amount)}
          </div>

          <div className="mb-6">
            <p className="text-gray-600">{product.description}</p>
          </div>

          <div className="mb-6">
            <p className="text-sm text-gray-500">カテゴリー: {product.category}</p>
            <p className="text-sm text-gray-500">SKU: {product.sku}</p>
            <p className="text-sm text-gray-500">在庫: {product.stock_quantity}個</p>
          </div>

          {/* 数量選択 */}
          <div className="mb-6">
            <label className="block text-sm font-medium mb-2">数量</label>
            <input
              type="number"
              min="1"
              max={product.stock_quantity}
              value={quantity}
              onChange={(e) => setQuantity(parseInt(e.target.value) || 1)}
              className="w-24 px-4 py-2 border rounded-lg"
            />
          </div>

          {/* アクションボタン */}
          <div className="flex gap-4">
            <button
              onClick={handleAddToCart}
              className="flex-1 bg-gray-200 text-gray-800 px-6 py-3 rounded-lg font-semibold hover:bg-gray-300 transition"
            >
              カートに追加
            </button>
            <button
              onClick={handleBuyNow}
              className="flex-1 bg-primary text-white px-6 py-3 rounded-lg font-semibold hover:bg-blue-600 transition"
            >
              今すぐ購入
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
