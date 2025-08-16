import { useEffect, useState } from 'react'
import { Package, Plus, RefreshCw, Pencil, Trash2, X } from 'lucide-react'
import { Category, Product } from '@/types'
import { apiService } from '@/services/api'
import { toast } from 'react-hot-toast'

const AdminCategoriesPage = () => {
	const [categories, setCategories] = useState<Category[]>([])
	const [products, setProducts] = useState<Product[]>([])
	const [loading, setLoading] = useState(true)
	const [saving, setSaving] = useState(false)
	const [showModal, setShowModal] = useState(false)
	const [editing, setEditing] = useState<Category | null>(null)
	const [form, setForm] = useState<{ name: string; description: string }>({ name: '', description: '' })

	useEffect(() => {
		fetchData()
	}, [])

	const fetchData = async () => {
		try {
			setLoading(true)
			const [cats, prods] = await Promise.all([
				apiService.getCategories(),
				apiService.getAllProducts()
			])
			setCategories(cats)
			setProducts(prods)
		} catch (e) {
			toast.error('Failed to load categories')
		} finally {
			setLoading(false)
		}
	}

	const usedCountByCategoryId = (): Record<number, number> => {
		const map: Record<number, number> = {}
		for (const p of products) {
			const cid = p.category_id
			if (cid) map[cid] = (map[cid] || 0) + 1
		}
		return map
	}

	const openCreate = () => {
		setEditing(null)
		setForm({ name: '', description: '' })
		setShowModal(true)
	}

	const openEdit = (cat: Category) => {
		setEditing(cat)
		setForm({ name: cat.name || '', description: cat.description || '' })
		setShowModal(true)
	}

	const handleSave = async () => {
		if (!form.name.trim()) {
			toast.error('Name is required')
			return
		}
		try {
			setSaving(true)
			if (editing) {
				await apiService.updateCategory(editing.id, { name: form.name.trim(), description: form.description })
				toast.success('Category updated')
			} else {
				await apiService.createCategory({ name: form.name.trim(), description: form.description })
				toast.success('Category created')
			}
			setShowModal(false)
			await fetchData()
		} catch (e) {
			toast.error('Failed to save category')
		} finally {
			setSaving(false)
		}
	}

	const handleDelete = async (cat: Category) => {
		const used = usedCountByCategoryId()[cat.id] || 0
		const msg = used > 0
			? `This category is used by ${used} product(s). Deleting it may orphan products. Continue?`
			: 'Delete this category?'
		if (!window.confirm(msg)) return
		try {
			await apiService.deleteCategory(cat.id)
			toast.success('Category deleted')
			await fetchData()
		} catch (e) {
			toast.error('Failed to delete category')
		}
	}

	return (
		<div className="space-y-6 p-6 bg-gray-50 dark:bg-gray-900 min-h-screen">
			<div className="flex items-center justify-between">
				<h1 className="text-3xl font-bold text-gray-900 dark:text-white">Categories</h1>
				<div className="flex gap-2">
					<button onClick={fetchData} className="bg-blue-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-blue-700 inline-flex items-center">
						<RefreshCw className="mr-2 h-4 w-4" /> Refresh
					</button>
					<button onClick={openCreate} className="bg-green-600 text-white px-4 py-2 rounded-lg font-medium hover:bg-green-700 inline-flex items-center">
						<Plus className="mr-2 h-4 w-4" /> Add Category
					</button>
				</div>
			</div>

			<div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
				{loading ? (
					<div className="p-6 text-gray-600 dark:text-gray-300">Loading...</div>
				) : categories.length === 0 ? (
					<div className="p-6 text-gray-600 dark:text-gray-300">No categories</div>
				) : (
					<ul className="divide-y divide-gray-200 dark:divide-gray-700">
						{categories.map(c => {
							const used = usedCountByCategoryId()[c.id] || 0
							return (
								<li key={c.id} className="p-4 flex items-center justify-between">
									<div className="flex items-center gap-3">
										<div className="bg-gray-100 dark:bg-gray-700 rounded-full p-2">
											<Package className="h-4 w-4" />
										</div>
										<div>
											<p className="font-medium text-gray-900 dark:text-white">{c.name}</p>
											<p className="text-sm text-gray-600 dark:text-gray-400">{c.description}</p>
											<div className="mt-1 text-xs inline-flex items-center px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300">
												{used} product(s)
											</div>
										</div>
									</div>
									<div className="flex items-center gap-2">
										<button onClick={() => openEdit(c)} className="px-3 py-1.5 rounded-md border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 inline-flex items-center">
											<Pencil className="h-4 w-4 mr-1" /> Edit
										</button>
										<button onClick={() => handleDelete(c)} className="px-3 py-1.5 rounded-md border border-red-300 text-red-700 hover:bg-red-50 inline-flex items-center">
											<Trash2 className="h-4 w-4 mr-1" /> Delete
										</button>
									</div>
								</li>
							)
						})}
					</ul>
				)}
			</div>

			{/* Modal */}
			{showModal && (
				<div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
					<div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg w-full max-w-md p-6">
						<div className="flex items-center justify-between mb-4">
							<h3 className="text-lg font-semibold text-gray-900 dark:text-white">{editing ? 'Edit Category' : 'Add Category'}</h3>
							<button onClick={() => setShowModal(false)} className="text-gray-500 hover:text-gray-700 dark:text-gray-300 dark:hover:text-white">
								<X className="h-5 w-5" />
							</button>
						</div>
						<div className="space-y-4">
							<div>
								<label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Name</label>
								<input
									type="text"
									value={form.name}
									onChange={(e) => setForm({ ...form, name: e.target.value })}
									className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
									placeholder="Category name"
								/>
							</div>
							<div>
								<label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Description</label>
								<textarea
									value={form.description}
									onChange={(e) => setForm({ ...form, description: e.target.value })}
									rows={3}
									className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
									placeholder="Describe this category"
								/>
							</div>
						</div>
						<div className="mt-6 flex justify-end gap-2">
							<button onClick={() => setShowModal(false)} className="px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200">Cancel</button>
							<button onClick={handleSave} disabled={saving} className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-60">
								{saving ? 'Saving...' : 'Save'}
							</button>
						</div>
					</div>
				</div>
			)}
		</div>
	)
}

export default AdminCategoriesPage
