import React, { useState } from 'react';
import { apiService } from '@/services/api';
import { toast } from 'react-hot-toast';

const ApiTestPage = () => {
  type TestResult = { success: boolean; data?: unknown; error?: string };
  const [testResults, setTestResults] = useState<Record<string, TestResult>>({});
  const [isLoading, setIsLoading] = useState<Record<string, boolean>>({});

  const runTest = async (testName: string, testFn: () => Promise<unknown>) => {
    setIsLoading(prev => ({ ...prev, [testName]: true }));
    try {
      const result = await testFn();
      setTestResults(prev => ({ ...prev, [testName]: { success: true, data: result } }));
      toast.success(`${testName} passed!`);
    } catch (error: unknown) {
      const errorMessage =
        typeof error === 'object' && error !== null && 'response' in error
          ? (error as { response?: { data?: { error?: string } }; message?: string }).response?.data?.error ||
            (error as { message?: string }).message ||
            'Unknown error'
          : 'Unknown error';
      setTestResults(prev => ({ ...prev, [testName]: { success: false, error: errorMessage } }));
      toast.error(`${testName} failed: ${errorMessage}`);
    } finally {
      setIsLoading(prev => ({ ...prev, [testName]: false }));
    }
  };

  const tests = [
    {
      name: 'Get Products',
      fn: () => apiService.getProducts()
    },
    {
      name: 'Get Categories',
      fn: () => apiService.getCategories()
    },
    {
      name: 'Get Orders',
      fn: () => apiService.getOrders()
    },
    {
      name: 'Get Conversations',
      fn: () => apiService.getConversations()
    }
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 mb-8">API Test Page</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Test Controls */}
        <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Run Tests</h2>
          <div className="space-y-3">
            {tests.map((test) => (
              <button
                key={test.name}
                onClick={() => runTest(test.name, test.fn)}
                disabled={isLoading[test.name]}
                className="w-full px-4 py-2 bg-blue-600 dark:bg-blue-500 text-white rounded-md hover:bg-blue-700 dark:hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                {isLoading[test.name] ? 'Running...' : `Test ${test.name}`}
              </button>
            ))}
          </div>
        </div>

        {/* Test Results */}
        <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Test Results</h2>
          <div className="space-y-4">
            {Object.entries(testResults).map(([testName, result]) => (
              <div key={testName} className="border border-gray-200 dark:border-gray-700 rounded-lg p-3 bg-gray-50 dark:bg-gray-800">
                <h3 className="font-medium text-gray-900 dark:text-gray-100 mb-2">{testName}</h3>
                {result.success ? (
                  <div className="text-green-600 dark:text-green-400">
                    <p>✅ Passed</p>
                    <pre className="text-xs mt-2 bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 p-2 rounded overflow-auto">
                      {JSON.stringify(result.data, null, 2)}
                    </pre>
                  </div>
                ) : (
                  <div className="text-red-600 dark:text-red-400">
                    <p>❌ Failed</p>
                    <p className="text-sm">{result.error}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ApiTestPage;
