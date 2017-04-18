// DCELVertex.cpp: implementation of the DCELVertex class.
//
//////////////////////////////////////////////////////////////////////

#include "stdafx.h"
#include "DCELVertex.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

//////////////////////////////////////////////////////////////////////
// Construction/Destruction
//////////////////////////////////////////////////////////////////////

DCELVertex::DCELVertex(): leaving(NULL), auxData(NULL), globalPrev(NULL), globalNext(NULL)
{
}

DCELVertex::~DCELVertex()
{

}

DCELHalfEdge* DCELVertex::getEdgeTo(const DCELVertex* v) const
{
	DCELHalfEdge* rval = NULL;

	if (leaving) {
		if (leaving->twin->origin == v) {
			rval = leaving;
		} else {
			DCELHalfEdge* test = leaving->twin->next;
			while (rval == NULL && test != leaving) {
				if (test->twin->origin == v) {
					rval = test;
				} else {
					test = test->twin->next;
				}
			}
		}
	}

	return rval;
}
