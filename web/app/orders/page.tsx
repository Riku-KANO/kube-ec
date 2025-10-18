'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { orderApi } from '@/lib/api';
import { useAuthStore } from '@/lib/store';
import type { Order } from '@/lib/types';

export default function OrdersPage() {
  const router = useRouter();
  const { user } = useAuthStore();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) {
      router.push('/login');
      return;
    }
    loadOrders();
  }, [user]);

  const loadOrders = async () => {
    if (!user) return;

    try {
      const data = await orderApi.list(user.id);
      setOrders(data.orders || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const formatPrice = (amount: number) => {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(amount);
  };

  const getStatusText = (status: string) => {
    const statusMap: { [key: string]: string } = {
      ORDER_STATUS_PENDING: '処理中',
      ORDER_STATUS_CONFIRMED: '確認済み',
      ORDER_STATUS_PROCESSING: '発送準備中',
      ORDER_STATUS_SHIPPED: '発送済み',
      ORDER_STATUS_DELIVERED: '配達完了',
      ORDER_STATUS_CANCELLED: 'キャンセル',
    };
    return statusMap[status] || status;
  };

  const getStatusColor = (status: string) => {
    const colorMap: { [key: string]: string } = {
      ORDER_STATUS_PENDING: 'bg-yellow-100 text-yellow-800',
      ORDER_STATUS_CONFIRMED: 'bg-blue-100 text-blue-800',
      ORDER_STATUS_PROCESSING: 'bg-purple-100 text-purple-800',
      ORDER_STATUS_SHIPPED: 'bg-green-100 text-green-800',
      ORDER_STATUS_DELIVERED: 'bg-gray-100 text-gray-800',
      ORDER_STATUS_CANCELLED: 'bg-red-100 text-red-800',
    };
    return colorMap[status] || 'bg-gray-100 text-gray-800';
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">読み込み中...</div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">注文履歴</h1>

      {orders.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 mb-4">まだ注文がありません</p>
          <button
            onClick={() => router.push('/products')}
            className="bg-primary text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
          >
            商品を見る
          </button>
        </div>
      ) : (
        <div className="space-y-6">
          {orders.map((order) => (
            <div key={order.id} className="border rounded-lg p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="font-semibold">注文番号: {order.id}</h3>
                  <p className="text-sm text-gray-600">
                    注文日: {new Date().toLocaleDateString('ja-JP')}
                  </p>
                </div>
                <span
                  className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                    order.status
                  )}`}
                >
                  {getStatusText(order.status)}
                </span>
              </div>

              <div className="space-y-2 mb-4">
                {order.items.map((item, index) => (
                  <div key={index} className="flex justify-between">
                    <span className="text-gray-600">
                      {item.product_name} × {item.quantity}
                    </span>
                    <span>{formatPrice(item.subtotal.amount)}</span>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4">
                <div className="flex justify-between font-bold">
                  <span>合計金額</span>
                  <span className="text-primary">{formatPrice(order.total_amount.amount)}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
