#include "stdafx.h"
#include "DCELMesh.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

DCELMesh::DCELMesh()
{
	faceList = NULL;
	vertexList = NULL;
	halfEdgeList = NULL;

	numFaces = 0;
	numHalfEdges = 0;
	numVertices = 0;
	min.x = min.y = min.z = -1.0;
	max.x = max.y = max.z = 1.0;
	vertexTotal.zero();
}

DCELMesh::~DCELMesh()
{
	clear();
}

void DCELMesh::clear()
{
	DCELFace* walkerF = faceList;
	DCELFace* tempF = NULL;
	while (walkerF) {
		tempF = walkerF->globalNext;
		delete walkerF;
		walkerF = tempF;
	}
	faceList = NULL;
	DCELVertex* walkerV = vertexList;
	DCELVertex* tempV = NULL;
	while (walkerV) {
		tempV = walkerV->globalNext;
		delete walkerV;
		walkerV = tempV;
	}
	vertexList = NULL;
	DCELHalfEdge* walkerE = halfEdgeList;
	DCELHalfEdge* tempE = NULL;
	while (walkerE) {
		tempE = walkerE->globalNext;
		delete walkerE;
		walkerE = tempE;
	}
	halfEdgeList = NULL;

	numFaces = 0;
	numHalfEdges = 0;
	numVertices = 0;
	min.x = min.y = min.z = -1.0;
	max.x = max.y = max.z = 1.0;
	vertexTotal.zero();
}

bool DCELMesh::isEmpty() const
{
	return ((vertexList == NULL) && (faceList == NULL) && (halfEdgeList == NULL));
}

void DCELMesh::insert(DCELVertex* v)
{
	if (v) {
		if (vertexList) {
			v->globalNext = vertexList;
			vertexList->globalPrev = v;
			vertexList = v;
		} else {
			vertexList = v;
		}
		numVertices++;
	}
}

void DCELMesh::insert(DCELFace* f)
{
	if (f) {
		if (faceList) {
			f->globalNext = faceList;
			faceList->globalPrev = f;
			faceList = f;
		} else {
			faceList = f;
		}
		numFaces++;
	}
}

void DCELMesh::insert(DCELHalfEdge* e)
{
	if (e) {
		if (halfEdgeList) {
			e->globalNext = halfEdgeList;
			halfEdgeList->globalPrev = e;
			halfEdgeList = e;
		} else {
			halfEdgeList = e;
		}
		numHalfEdges++;
	}
}

void DCELMesh::remove(DCELVertex* v)
{
	if (v) {
		if (vertexList == v) {
			vertexList = vertexList->globalNext;
			if (vertexList) {
				vertexList->globalPrev = NULL;
			}
		} else {
			v->globalPrev->globalNext = v->globalNext;
			if (v->globalNext) {
				v->globalNext->globalPrev = v->globalPrev;
			}
		}
		v->globalNext = NULL;
		v->globalPrev = NULL;
		numVertices--;
	}
}

void DCELMesh::remove(DCELFace* f)
{
	if (f) {
		if (faceList == f) {
			faceList = faceList->globalNext;
			if (faceList) {
				faceList->globalPrev = NULL;
			}
		} else {
			f->globalPrev->globalNext = f->globalNext;
			if (f->globalNext) {
				f->globalNext->globalPrev = f->globalPrev;
			}
		}
		f->globalNext = NULL;
		f->globalPrev = NULL;
		numFaces--;
	}
}

void DCELMesh::remove(DCELHalfEdge* e)
{
	if (e) {
		if (halfEdgeList == e) {
			halfEdgeList = halfEdgeList->globalNext;
			if (halfEdgeList) {
				halfEdgeList->globalPrev = NULL;
			}
		} else {
			e->globalPrev->globalNext = e->globalNext;
			if (e->globalNext) {
				e->globalNext->globalPrev = e->globalPrev;
			}
		}
		e->globalNext = NULL;
		e->globalPrev = NULL;
		numHalfEdges--;
	}
}

void DCELMesh::updateFaceNormals()
{
	for (DCELFace* walker = faceList; walker; advance(walker)) {
		walker->updateNormal();
	}
}

void DCELMesh::updateVertexNormals()
{
	DCELVertex* walkerV;
	for (walkerV = vertexList; walkerV; advance(walkerV)) {
		walkerV->normal.zero();
	}
	for (DCELFace* walkerF = faceList; walkerF ;advance(walkerF)) {
		walkerF->updateVertexNormals();
	}
	for (walkerV = vertexList; walkerV; advance(walkerV)) {
		walkerV->normal.normalize();
	}
}

void DCELMesh::updateEdgeBits()
{
	for (DCELHalfEdge* walkerE = halfEdgeList; walkerE; advance(walkerE)) {
		if (walkerE->face == &infiniteFace ||
			walkerE->twin->face == &infiniteFace) {
			walkerE->setMask(DCELHalfEdge::DCEL_EDGE_BOUNDARY_BIT, true);
		} else {
			walkerE->setMask(DCELHalfEdge::DCEL_EDGE_BOUNDARY_BIT, false);
		}
	}
}

void DCELMesh::updateStatistics()
{
	vertexTotal.zero();

	DCELVertex* walkerV = vertexList;
	if (walkerV) {
		min = walkerV->coords;
		max = walkerV->coords;
		vertexTotal.translateBy(walkerV->coords);
		walkerV = walkerV->globalNext;
	}
	while (walkerV) {
		if (walkerV->coords.x < min.x) {
			min.x = walkerV->coords.x;
		} else if (walkerV->coords.x > max.x) {
			max.x = walkerV->coords.x;
		}
		if (walkerV->coords.y < min.y) {
			min.y = walkerV->coords.y;
		} else if (walkerV->coords.y > max.y) {
			max.y = walkerV->coords.y;
		}
		if (walkerV->coords.z < min.z) {
			min.z = walkerV->coords.z;
		} else if (walkerV->coords.z > max.z) {
			max.z = walkerV->coords.z;
		}
		vertexTotal.translateBy(walkerV->coords);
		advance(walkerV);
	}
}

Vector DCELMesh::getCentroid() const
{
	return vertexTotal * (1.0 / (double)numVertices);
}

void DCELMesh::updateAll()
{
	updateFaceNormals();
	updateVertexNormals();
	updateEdgeBits();
	updateStatistics();
}

int DCELMesh::getNumTriangles()
{
	int rval = 0;
	
	for (DCELFace* walkerF = faceList; walkerF; advance(walkerF)) {
		if (walkerF->edge->next->next->next == walkerF->edge) {
			rval++;
		}
	}
	return rval;
}

int DCELMesh::getNumQuads()
{
	int rval = 0;
	
	for (DCELFace* walkerF = faceList; walkerF; advance(walkerF)) {
		if (walkerF->edge->next->next->next->next == walkerF->edge) {
			rval++;
		}
	}
	return rval;
}

void DCELMesh::setHalfEdgeMasks(unsigned int mask, bool value)
{
	for (DCELHalfEdge* walkerE = halfEdgeList; walkerE; advance(walkerE)) {
		walkerE->setMask(mask, value);
	}
}

