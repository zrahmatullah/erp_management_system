export interface Category { id: string; name: string; }
export interface Product { id: string; name: string; categoryId: string; price: number; }
export interface ProductVariant { id: string; productId: string; name: string; additionalPrice: number; }
export interface ProductImage { id: string; url: string; }
export interface Recipe { id: string; productId: string; }
export interface RecipeItem { id: string; recipeId: string; inventoryItemId: string; quantity: number; }
export interface ProductCreateRequest { name: string; price: number; categoryId: string; }
export interface ProductFilters { search?: string; categoryId?: string; }
