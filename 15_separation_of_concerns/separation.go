package separation

import "fmt"

// ══════════════════════════════════════════════════
// Separation of Concerns (SoC)
// Setiap layer punya tanggung jawab sendiri.
// ══════════════════════════════════════════════════

// ── MODEL LAYER ──
// Hanya mendefinisikan data structure

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

// ── REPOSITORY LAYER ──
// Bertanggung jawab atas data access (CRUD)

type ProductRepository interface {
	FindAll() []Product
	FindByID(id int) (Product, bool)
	Save(p Product)
}

type InMemoryProductRepo struct {
	products []Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		products: []Product{
			{ID: 1, Name: "Laptop", Price: 15000000, Stock: 10},
			{ID: 2, Name: "Mouse", Price: 250000, Stock: 50},
			{ID: 3, Name: "Keyboard", Price: 500000, Stock: 30},
		},
	}
}

func (r *InMemoryProductRepo) FindAll() []Product {
	return r.products
}

func (r *InMemoryProductRepo) FindByID(id int) (Product, bool) {
	for _, p := range r.products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

func (r *InMemoryProductRepo) Save(p Product) {
	r.products = append(r.products, p)
}

// ── SERVICE LAYER ──
// Business logic — tidak peduli data datang dari mana

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() []Product {
	return s.repo.FindAll()
}

func (s *ProductService) GetProduct(id int) (Product, error) {
	p, found := s.repo.FindByID(id)
	if !found {
		return Product{}, fmt.Errorf("product with ID %d not found", id)
	}
	return p, nil
}

func (s *ProductService) IsExpensive(p Product) bool {
	return p.Price > 1000000
}

// ── PRESENTER / VIEW LAYER ──
// Hanya menangani tampilan output

func PrintProduct(p Product, expensive bool) {
	label := ""
	if expensive {
		label = " 💰 PREMIUM"
	}
	fmt.Printf("  [%d] %s — Rp%.0f (stock: %d)%s\n", p.ID, p.Name, p.Price, p.Stock, label)
}

func Run() {
	fmt.Println("=== Separation of Concerns ===")

	// Wire up layers
	repo := NewInMemoryProductRepo()
	service := NewProductService(repo)

	// ── List all products ──
	fmt.Println("\nSemua Produk:")
	products := service.GetAllProducts()
	for _, p := range products {
		PrintProduct(p, service.IsExpensive(p))
	}

	// ── Find by ID ──
	fmt.Println("\nCari Product ID 2:")
	p, err := service.GetProduct(2)
	if err != nil {
		fmt.Println("  Error:", err)
	} else {
		PrintProduct(p, service.IsExpensive(p))
	}

	// ── Not found ──
	fmt.Println("\nCari Product ID 99:")
	_, err = service.GetProduct(99)
	if err != nil {
		fmt.Println("  Error:", err)
	}

	// ── Architecture overview ──
	fmt.Println("\n=== Layer Architecture ===")
	fmt.Println("  Model      → Definisi data (struct)")
	fmt.Println("  Repository → Data access (CRUD)")
	fmt.Println("  Service    → Business logic")
	fmt.Println("  Presenter  → Output / tampilan")
	fmt.Println("\n  Setiap layer hanya berkomunikasi melalui interface.")
	fmt.Println("  Repo bisa diganti (DB, File, API) tanpa ubah Service.")
}
