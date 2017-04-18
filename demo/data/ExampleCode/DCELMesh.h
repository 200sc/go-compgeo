#if !defined(AFX_DCELMESH_H__92DA9CFE_BA7F_4A5C_8A08_C10711268F28__INCLUDED_)
#define AFX_DCELMESH_H__92DA9CFE_BA7F_4A5C_8A08_C10711268F28__INCLUDED_

#if _MSC_VER > 1000
#pragma once
#endif // _MSC_VER > 1000

#include "DCELFace.h"
#include "DCELVertex.h"
#include "DCELHalfEdge.h"

/*
 * DCELMesh class. Part of an example DCEL implementation
 * Webpage: http://www.holmes3d.net/graphics/dcel/
 * Author: Ryan Holmes
 * E-mail: ryan <at> holmes3d <dot> net
 * Usage: Use freely. Please cite the website as the source if you
 * use it substantially unchanged. Please leave this documentation
 * in the code.
 */

class DCELMesh  
{
public:
	DCELMesh();
	~DCELMesh();

	// Simple iteration interface. Supports forward traversal only
	DCELFace* firstFace() { return faceList; };
	DCELFace* next(DCELFace* f) { return (f != NULL) ? f->globalNext : NULL; };
	void advance(DCELFace* &f) { f = (f != NULL) ? f->globalNext : NULL; };

	DCELVertex* firstVertex() { return vertexList; };
	DCELVertex* next(DCELVertex* v) { return (v != NULL) ? v->globalNext : NULL; };
	void advance(DCELVertex* &v) { v = (v != NULL) ? v->globalNext : NULL; };

	DCELHalfEdge* firstHalfEdge() { return halfEdgeList; };
	DCELHalfEdge* next(DCELHalfEdge* e) { return (e != NULL) ? e->globalNext : NULL; };
	void advance(DCELHalfEdge* &e) { e = (e != NULL) ? e->globalNext : NULL; };

	// Deletes all objects in the mesh. Does not delete auxData objects.
	void clear();
	// True iff the mesh contains no objects
	bool isEmpty() const;

	// Inserts the object at the head of the appropriate list. This means that insertion
	// is a safe operation during processing of all objects of a particular type, because
	// they will be inserted before the current iterator position.
	void insert(DCELVertex* v);
	void insert(DCELFace* f);
	void insert(DCELHalfEdge* e);

	// Removes from the mesh holder, but does not delete or disconnect. Caller is responsible
	// for correct usage. Removing an object that is not in the mesh is an unsafe operation.
	void remove(DCELVertex* v);
	void remove(DCELFace* f);
	void remove(DCELHalfEdge* e);

	// Returns the one and only infiniteFace pointer. This never changes during
	// the lifetime of the mesh.
	DCELFace* getInfiniteFace() { return &infiniteFace; };

	// Calculates a normal for each face based on its geometry
	void updateFaceNormals();
	// Calculates a normal for each vertex based on the normals of the faces around it.
	// Requires updateFaceNormals() call before, or manual setting of face normals.
	void updateVertexNormals();
	// Sets purely geometric edge bits. Does not affect viewing or view-specific bits
	void updateEdgeBits();
	// Calculates bounding box and internal statistics
	void updateStatistics();
	// Shorthand for calling the above four update functions
	void updateAll();

	// Return current counts of member objects
	int getNumFaces() const { return numFaces; };
	int getNumVertices() const { return numVertices; };
	int getNumHalfEdges() const { return numHalfEdges; };

	// Calculate and return values
	int getNumTriangles();
	int getNumQuads();

	// Helper function to set or clear a particular mask on all HalfEdges
	void setHalfEdgeMasks(unsigned int mask, bool value);

	// Only valid after an updateStatistics or updateAll call.
	// Returns center of mass of object, assuming all vertices have equal mass
	Vector getCentroid() const;
	void loadBoundingBox(Vector &minPoint, Vector &maxPoint) const { minPoint = min; maxPoint = max; };

private:
	DCELFace* faceList;
	DCELVertex* vertexList;
	DCELHalfEdge* halfEdgeList;

	int numFaces;
	int numVertices;
	int numHalfEdges;

	Vector min;
	Vector max;
	Vector vertexTotal;

	DCELFace infiniteFace;
};

#endif // !defined(AFX_DCELMESH_H__92DA9CFE_BA7F_4A5C_8A08_C10711268F28__INCLUDED_)
