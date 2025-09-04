import React, { createContext, useContext, useState, type ReactNode } from 'react';
import { ToastContainer, type ToastData } from '../components/Toast';

interface ToastContextType {
  showToast: (message: string, type?: 'info' | 'success' | 'warning' | 'error', duration?: number) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = () => {
  const context = useContext(ToastContext);
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

interface ToastProviderProps {
  children: ReactNode;
}

export const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const showToast = (
    message: string, 
    type: 'info' | 'success' | 'warning' | 'error' = 'info',
    duration: number = 3000
  ) => {
    const id = Date.now().toString() + Math.random().toString(36).substr(2, 9);
    const newToast: ToastData = {
      id,
      message,
      type,
      duration,
    };

    setToasts(prev => [...prev, newToast]);
  };

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id));
  };

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <ToastContainer toasts={toasts} onRemoveToast={removeToast} />
    </ToastContext.Provider>
  );
};
