// DCELHalfEdge.cpp: implementation of the DCELHalfEdge class.
//
//////////////////////////////////////////////////////////////////////

#include "stdafx.h"
#include "DCELHalfEdge.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

//////////////////////////////////////////////////////////////////////
// Construction/Destruction
//////////////////////////////////////////////////////////////////////

DCELHalfEdge::DCELHalfEdge() :
twin(NULL), next(NULL), face(NULL), origin(NULL), auxData(NULL), displayBits(0),
globalPrev(NULL), globalNext(NULL)
{
}

DCELHalfEdge::~DCELHalfEdge()
{

}

DCELHalfEdge* DCELHalfEdge::getPrev()
{
	DCELHalfEdge* rval = twin->next->twin;
	
	while (rval->next != this) {
		rval = rval->next->twin;
	}

	return rval;
}
