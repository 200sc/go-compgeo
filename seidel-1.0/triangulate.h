

/* Node attributes for every node in the query structure */

typedef struct {
  int vnum;
  int next;			/* Circularly linked list  */
  int prev;			/* describing the monotone */
  int marked;			/* polygon */
} monchain_t;			


typedef struct {
  point_t pt;
  int vnext[4];			/* next vertices for the 4 chains */
  int vpos[4];			/* position of v in the 4 chains */
  int nextfree;
} vertexchain_t;


/* Node types */

#define T_X     1
#define T_Y     2
#define T_SINK  3

#define FIRSTPT 1		/* checking whether pt. is inserted */ 
#define LASTPT  2


#define S_LEFT 1		/* for merge-direction */
#define S_RIGHT 2


#define ST_VALID 1		/* for trapezium state */
#define ST_INVALID 2


#define SP_SIMPLE_LRUP 1	/* for splitting trapezoids */
#define SP_SIMPLE_LRDN 2
#define SP_2UP_2DN     3
#define SP_2UP_LEFT    4
#define SP_2UP_RIGHT   5
#define SP_2DN_LEFT    6
#define SP_2DN_RIGHT   7
#define SP_NOSPLIT    -1	

#define TR_FROM_UP 1		/* for traverse-direction */
#define TR_FROM_DN 2

#define TRI_LHS 1
#define TRI_RHS 2

/* Functions */

extern int monotonate_trapezoids(int);
extern int triangulate_monotone_polygons(int, int, int (*)[3]);

extern int locate_endpoint(point_t *, point_t *, int);
extern int construct_trapezoids(int);

extern int choose_segment(void);
extern int read_segments(char *, int *);

