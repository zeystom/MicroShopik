import React from 'react'
import { Link } from 'react-router-dom'
import { 
  Gamepad2, 
  Monitor, 
  Smartphone, 
  Tv, 
  Headphones,
  Shield,
  Clock,
  Users
} from 'lucide-react'

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear()

  const categories = [
    { name: 'Gaming Accounts', icon: Gamepad2, href: '/products?category=gaming' },
    { name: 'Streaming Services', icon: Monitor, href: '/products?category=streaming' },
    { name: 'Mobile Apps', icon: Smartphone, href: '/products?category=mobile' },
    { name: 'Software Licenses', icon: Tv, href: '/products?category=software' },
    { name: 'Digital Music', icon: Headphones, href: '/products?category=music' }
  ]

  const quickLinks = [
    { name: 'About Us', href: '/about' },
    { name: 'How It Works', href: '/how-it-works' },
    { name: 'Seller Guidelines', href: '/seller-guidelines' },
    { name: 'Buyer Protection', href: '/buyer-protection' },
    { name: 'Terms of Service', href: '/terms' },
    { name: 'Privacy Policy', href: '/privacy' }
  ]

  const supportLinks = [
    { name: 'Help Center', href: '/help' },
    { name: 'Contact Support', href: '/contact' },
    { name: 'Dispute Resolution', href: '/disputes' },
    { name: 'FAQ', href: '/faq' }
  ]

  const features = [
    {
      icon: Shield,
      title: 'Secure Transactions',
      description: 'End-to-end encryption and secure payment processing'
    },
    {
      icon: Clock,
      title: 'Instant Delivery',
      description: 'Get your digital goods immediately after purchase'
    },
    {
      icon: Users,
      title: 'Verified Sellers',
      description: 'All sellers are verified and rated by our community'
    }
  ]

  return (
    <footer className="bg-slate-900 text-slate-100">
      <div className="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8 py-8 sm:py-12">
        {/* Main Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 sm:gap-8 mb-8 sm:mb-12">
          {/* Company Info */}
          <div className="lg:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">MS</span>
              </div>
              <span className="text-xl font-bold">MicroShopik</span>
            </div>
            <p className="text-slate-300 mb-4 sm:mb-6 max-w-md text-sm sm:text-base">
              Your trusted marketplace for premium digital goods. Buy and sell gaming accounts, 
              streaming services, software licenses, and more with confidence.
            </p>
            
            {/* Features */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4">
              {features.map((feature) => {
                const Icon = feature.icon
                return (
                  <div key={feature.title} className="flex items-start space-x-3">
                    <Icon className="w-5 h-5 text-blue-400 mt-0.5 flex-shrink-0" />
                    <div>
                      <h4 className="font-medium text-white text-sm">{feature.title}</h4>
                      <p className="text-gray-400 text-xs sm:text-sm">{feature.description}</p>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>

          {/* Categories */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Categories</h3>
            <ul className="space-y-2">
              {categories.map((category) => {
                const Icon = category.icon
                return (
                  <li key={category.name}>
                    <Link
                      to={category.href}
                      className="flex items-center space-x-2 text-slate-300 hover:text-slate-100 transition-colors"
                    >
                      <Icon className="w-4 h-4" />
                      <span>{category.name}</span>
                    </Link>
                  </li>
                )
              })}
            </ul>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-lg font-semibold mb-4">Quick Links</h3>
            <ul className="space-y-2">
              {quickLinks.map((link) => (
                <li key={link.name}>
                  <Link
                    to={link.href}
                    className="text-slate-300 hover:text-slate-100 transition-colors"
                  >
                    {link.name}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Support Section */}
        <div className="border-t border-slate-800 pt-8 mb-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div>
              <h3 className="text-lg font-semibold mb-4">Support</h3>
              <ul className="space-y-2">
                {supportLinks.map((link) => (
                  <li key={link.name}>
                    <Link
                      to={link.href}
                      className="text-slate-300 hover:text-slate-100 transition-colors"
                    >
                      {link.name}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Newsletter</h3>
              <p className="text-slate-300 mb-4">
                Stay updated with the latest digital goods and exclusive offers.
              </p>
              <div className="flex space-x-2">
                <input
                  type="email"
                  placeholder="Enter your email"
                  className="flex-1 px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-slate-100 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
                <button className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
                  Subscribe
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom Footer */}
        <div className="border-t border-slate-800 pt-8">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <div className="text-slate-400 text-sm">
              © {currentYear} MicroShopik. All rights reserved.
            </div>
            
            <div className="flex items-center space-x-6 text-sm text-slate-400">
              <span>Trusted by 10,000+ users</span>
              <span>•</span>
              <span>24/7 Support</span>
              <span>•</span>
              <span>Secure Platform</span>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}

export default Footer

