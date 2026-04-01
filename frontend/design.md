# Design System: Bauhaus Modernist

### 1. Overview & Creative North Star
**Creative North Star: "Form as Function"**
Bauhaus Modernist is a design system that rejects the "softness" of modern SaaS interfaces in favor of architectural rigor, primary colors, and heavy structural lines. It is inspired by the early 20th-century Bauhaus movement—emphasizing that design is utility and information is structure.

The system breaks the "template" look through:
- **Neo-Brutalist Anchoring:** Using heavy black borders (4px) to define the perimeter of the experience.
- **Intentional Asymmetry:** Utilizing bento-style grids where weight is distributed through color blocks rather than centered alignments.
- **Kinetic Interaction:** High-contrast hover states and "physical" active states where elements shift 1-2px to simulate mechanical depression.

### 2. Colors
The palette is built on the triad of Red, Blue, and Yellow, anchored by absolute Black and White.

- **The "Heavy Line" Rule:** Unlike systems that prohibit borders, Bauhaus Modernist mandates them. Use 4px solid black (`#000000`) or white (`#ffffff`) borders for all major structural containers (Header, Sidebar, Main Section).
- **Surface Hierarchy:**
- `Surface`: Pure White (#ffffff) for maximum contrast.
- `Surface Container Low`: A soft mint-tint (#F1FAEE) used exclusively for active navigation states.
- `Surface Container`: Light grey (#eeeeee) for row-level hover states.
- **Primary Roles:** Use `Primary` (Red) for critical alerts and CTAs. Use `Secondary` (Blue) for user-specific data. Use `Tertiary` (Yellow/Gold) for active processes.
- **Signature Textures:** Avoid gradients. Use solid blocks of primary colors to define sections.

### 3. Typography
The system uses a high-contrast pairing: **Space Grotesk** for structural elements and **Work Sans** for utilitarian data.

**Typography Scale:**
- **Display (8xl):** 4.5rem (72px) - Used for major greetings or section heroes. Always uppercase, black-weight, tracking-tighter.
- **Headline (text-3xl):** 1.875rem (30px) - Used for card titles and section headers.
- **Body (0.875rem - 1rem):** Used for descriptions and feed items.
- **Label (10px - 12px):** Used for metadata, always uppercase with wide tracking (0.2em to 0.4em) to create an "architectural blueprint" feel.

The typographic hierarchy conveys identity through "Scale Extremes"—pairing very large headlines with very small, wide-tracked labels.

### 4. Elevation & Depth
Traditional depth (z-index shadows) is replaced by **Hard-Shadow Layering**.

- **The Layering Principle:** Depth is created by "Stacking" containers with offset black borders.
- **Ambient Shadows:** Standard CSS box-shadows are prohibited. Instead, use "Hard Shadows": `shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]`.
- **The "Tactile Press":** On click/active, elements should lose their hard shadow and translate: `active:translate-x-1 active:translate-y-1 active:shadow-none`.
- **Glassmorphism:** Prohibited. Everything must be opaque and definitive.

### 5. Components
- **Buttons:** Rectangular, 0px radius, 2px or 4px black borders. High-contrast hover (e.g., White to Red).
- **Cards:** Defined by 4px black perimeters. No rounded corners. Internal padding is generous (p-8).
- **Navigation:** Vertical sidebar with high-contrast active states using `Surface Container Low` and a 2px border.
- **Bento Stats:** Blocks of solid primary colors with white text for high-impact data visualization.
- **FAB:** A square red box with a heavy 4px border and an 8px hard shadow.

### 6. Do's and Don'ts
**Do:**
- Use absolute black (#000000) for all structural borders.
- Ensure all headlines are uppercase with tight tracking.
- Use 0px or very minimal (4px) border-radius to maintain geometric purity.
- Lean into asymmetrical layouts.

**Don't:**
- Use subtle 1px grey borders.
- Use soft, diffused shadows or blurs.
- Use pastel colors for primary actions.
- Mix rounded corners with sharp corners; stick to the geometric 0px rule.