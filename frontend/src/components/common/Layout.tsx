import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import Footer from './Footer';

const Layout = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Header />
      <main className="p-4 md:p-6">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
};

export default Layout;

