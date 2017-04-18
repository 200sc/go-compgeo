#include "stdafx.h"
#include "DCELFace.h"
#include "DCELVertex.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

DCELFace::DCELFace() : edge(NULL), auxData(NULL), globalPrev(NULL), globalNext(NULL)
{
}

DCELFace::~DCELFace()
{

}

void DCELFace::updateNormal()
{
	normal.zero();
	Vector backEdge;
	Vector forwardEdge;
	DCELHalfEdge* walker = edge;
	if (walker) {
		forwardEdge = walker->twin->origin->coords - walker->origin->coords;
		walker = walker->next;
	}
	while (walker != edge) {
		backEdge = forwardEdge * -1.0;
		forwardEdge = walker->twin->origin->coords - walker->origin->coords;
		normal.translateBy(forwardEdge.Cross(backEdge));
		walker = walker->next;
	}
	backEdge = forwardEdge * -1.0;
	forwardEdge = walker->twin->origin->coords - walker->origin->coords;
	normal.translateBy(forwardEdge.Cross(backEdge));
	normal.normalize();
}

void DCELFace::updateVertexNormals() const
{
	DCELHalfEdge* walker = edge;
	if (walker) {
		walker->origin->normal.translateBy(normal);
		walker = walker->next;
	}
	while (walker != edge) {
		walker->origin->normal.translateBy(normal);
		walker = walker->next;
	}
}

int DCELFace::getEdgeCount() const
{
	int rval = 0;
	if (edge) {
		DCELHalfEdge* walkerE = edge->next;
		rval = 1;
		while (walkerE != edge) {
			rval++;
			walkerE = walkerE->next;
		}
	}
	return rval;
}
