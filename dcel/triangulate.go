package dcel

/* Input specified as contours.
 * Outer contour must be anti-clockwise.
 * All inner contours must be clockwise.
 *
 * Every contour is specified by giving all its points in order. No
 * point shoud be repeated. i.e. if the outer contour is a square,
 * only the four distinct endpoints shopudl be specified in order.
 *
 * ncontours: #contours
 * cntr: An array describing the number of points in each
 *	 contour. Thus, cntr[i] = #points in the i'th contour.
 * vertices: Input array of vertices. Vertices for each contour
 *           immediately follow those for previous one. Array location
 *           vertices[0] must NOT be used (i.e. i/p starts from
 *           vertices[1] instead. The output triangles are
 *	     specified  w.r.t. the indices of these vertices.
 * triangles: Output array to hold triangles.
 *
 * Enough space must be allocated for all the arrays before calling
 * this routine
 */

// func (f *Face) Triangulate() *DCEL {}

// func (dc *DCEL) Triangulate() *DCEL {
//      int cntr[];
//      double (*vertices)[2];
// {

//   while (ccount < ncontours)
//     {
//       int j;
//       int first, last;

//       npoints = cntr[ccount];
//       first = i;
//       last = first + npoints - 1;
//       for (j = 0; j < npoints; j++, i++)
// 	{
// 	  seg[i].v0.x = vertices[i][0];
// 	  seg[i].v0.y = vertices[i][1];

// 	  if (i == last)
// 	    {
// 	      seg[i].next = first;
// 	      seg[i].prev = i-1;
// 	      seg[i-1].v1 = seg[i].v0;
// 	    }
// 	  else if (i == first)
// 	    {
// 	      seg[i].next = i+1;
// 	      seg[i].prev = last;
// 	      seg[last].v1 = seg[i].v0;
// 	    }
// 	  else
// 	    {
// 	      seg[i].prev = i-1;
// 	      seg[i].next = i+1;
// 	      seg[i-1].v1 = seg[i].v0;
// 	    }

// 	  seg[i].is_inserted = FALSE;
// 	}

//       ccount++;
//     }

//   genus = ncontours - 1;
//   n = i-1;

//   // Randomize_segements()
//   construct_trapezoids(n)
//   nmonpoly = monotonate_trapezoids(n)
//   triangulate_monotone_polygons(n, nmonpoly, triangles)
// }

type Trapezoid struct {
	Left, Right          *Edge
	MinY, MaxY           float64
	Valid                bool
	u0, u1, uSave, uSide int
	d0, d1               int
	sink                 *Node
}

type NodeType int

// NodeType const
const (
	T_X NodeType = iota
	T_Y
	T_SINK
)

type Node struct {
	e                   *Edge
	t                   *Trapezoid
	typ                 NodeType
	yval                Point
	parent, left, right *Node
}

type TrapEdge struct {
	*Edge
	r0, r1 *Node
}
