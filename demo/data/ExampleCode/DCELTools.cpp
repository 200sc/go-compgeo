// DCELTools.cpp: implementation of the DCELTools class.
//
//////////////////////////////////////////////////////////////////////

#include "stdafx.h"
#include "DCELTools.h"
#include "HalfEdgeList.h"

#include <fstream>
#include <string>

using namespace std;

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#define new DEBUG_NEW
#endif

//////////////////////////////////////////////////////////////////////
// Construction/Destruction
//////////////////////////////////////////////////////////////////////

DCELTools::DCELTools()
{

}

DCELTools::~DCELTools()
{

}

bool DCELTools::isConsistent(DCELHalfEdge* e)
{
	bool rval = true;

	if ((e->twin == NULL) ||
		(e->twin->twin != e) ||
		(e->origin == NULL) ||
		(e->face == NULL) ||
		(e->next == NULL) ||
		(e->next->face != e->face)) {
		rval = false;
	}

	return rval;
}

bool DCELTools::isConsistent(DCELFace* f)
{
	bool rval = true;

	if ((f->edge == NULL) ||
		(f->edge->face != f)) {
		rval = false;
	}

	return rval;
}

bool DCELTools::isConsistent(DCELVertex* v)
{
	bool rval = true;

	if ((v->leaving == NULL) ||
		(v->leaving->origin != v)) {
		rval = false;
	}

	return rval;
}

bool DCELTools::isConsistent(DCELMesh* m)
{
	bool rval = true;

	DCELHalfEdge* eWalker = m->firstHalfEdge();
	while (rval && eWalker) {
		if (!isConsistent(eWalker)) {
			rval = false;
		} else {
			m->advance(eWalker);
		}
	}
	DCELVertex* vWalker = m->firstVertex();
	while (rval && vWalker) {
		if (!isConsistent(vWalker)) {
			rval = false;
		} else {
			m->advance(vWalker);
		}
	}
	DCELFace* fWalker = m->firstFace();
	while (rval && fWalker) {
		if (!isConsistent(fWalker)) {
			rval = false;
		} else {
			m->advance(fWalker);
		}
	}

	return rval;
}

bool DCELTools::storeToOFF(LPCTSTR filename, DCELMesh* m)
{
	bool rval = false;

	ofstream ofs(filename);

	if (ofs.is_open()) {
		// Create a temporary index array to associate with vertices through
		// their auxData pointers. These will be used as ids during the write
		// process. The existing pointers are stored in case they were in use
		// before the save, and are restored at the end of this function.
		unsigned int pointCounter = m->getNumVertices();
		void** pointerBuffer = new void*[pointCounter];
		unsigned int* indexBuffer = new unsigned int[pointCounter];
		pointCounter = 0;
		DCELVertex* walkerV;
		for (walkerV = m->firstVertex(); walkerV; m->advance(walkerV)) {
			pointerBuffer[pointCounter] = walkerV->auxData;
			indexBuffer[pointCounter] = pointCounter;
			walkerV->auxData = &(indexBuffer[pointCounter]);
			pointCounter++;
		}
		int numFaces = m->getNumFaces();
		int numEdges = m->getNumHalfEdges() / 2;

		ofs << "OFF" << endl;
		ofs << pointCounter << " " << numFaces << " " << numEdges << endl;
		for (walkerV = m->firstVertex(); walkerV; m->advance(walkerV)) {
			ofs << walkerV->coords.x << " " << walkerV->coords.y << " " << walkerV->coords.z;
			ofs << endl;
		}
		// For each poly, count the number of edges it has,
		// then go around and actually output the ids of the
		// vertices.
		int polySize;
		
		DCELHalfEdge* walkerE;
		for (DCELFace* walkerF = m->firstFace(); walkerF; m->advance(walkerF)) {
			polySize = 0;
			walkerE = walkerF->edge;
			polySize++;
			walkerE = walkerE->next;
			while (walkerE != walkerF->edge) {
				polySize++;
				walkerE = walkerE->next;
			}
			ofs << polySize << " " << *((unsigned int*)(walkerE->origin->auxData));
			walkerE = walkerE->next;
			while (walkerE != walkerF->edge) {
				ofs << " " << *((unsigned int*)(walkerE->origin->auxData));
				walkerE = walkerE->next;
			}
			ofs << endl;
		}

		// Set the auxData pointers back.
		pointCounter = 0;
		for (walkerV = m->firstVertex(); walkerV ; m->advance(walkerV)) {
			walkerV->auxData = pointerBuffer[pointCounter];
			pointCounter++;
		}
		delete[] indexBuffer;
		delete[] pointerBuffer;
		ofs.close();
		rval = true;
	}

	return rval;
}

bool DCELTools::loadFromOFF(LPCTSTR filename, DCELMesh* m)
{
	bool rval = false;

	ifstream ifs(filename);

	m->clear();

	if (ifs.is_open()) {

		string inputLine;
		bool isManifold = true;

		::getline(ifs, inputLine);
		if (inputLine == "OFF") {
			int numVertices = 0;
			int numFaces = 0;
			int numEdges;

			ifs >> numVertices >> numFaces >> numEdges; // numEdges is present but ignored in this format

			if (numVertices > 0 && numFaces > 0) {
				DCELVertex* tempVertex = NULL;
				DCELHalfEdge* tempHalfEdge = NULL;
				DCELFace* tempFace = NULL;

				// Temporary array to give us indexed access to the vertices during the build
				DCELVertex** vertices = new DCELVertex*[numVertices];

				for (int i = 0; i < numVertices; i++) {
					tempVertex = new DCELVertex();
					ifs >> tempVertex->coords.x >> tempVertex->coords.y >> tempVertex->coords.z;
					vertices[i] = tempVertex;
					m->insert(tempVertex);
				}

				int vIndex;
				HalfEdgeList* pList;
				for (int i = 0; i < numFaces; i++) {
					ifs >> numEdges; // Number of edges/vertices in this polygon
					tempFace = new DCELFace();
					m->insert(tempFace);
					tempHalfEdge = new DCELHalfEdge();
					m->insert(tempHalfEdge);
					tempFace->edge = tempHalfEdge;
					tempHalfEdge->face = tempFace;
					ifs >> vIndex;
					tempHalfEdge->origin = vertices[vIndex];
					vertices[vIndex]->leaving = tempHalfEdge;
					pList = (HalfEdgeList*)vertices[vIndex]->auxData;
					HalfEdgeList::addToList(pList, tempHalfEdge);
					vertices[vIndex]->auxData = pList;
					for (int j = 1; j < numEdges; j++) {
						tempHalfEdge->next = new DCELHalfEdge();
						tempHalfEdge = tempHalfEdge->next;
						m->insert(tempHalfEdge);
						tempHalfEdge->face = tempFace;
						ifs >> vIndex;
						tempHalfEdge->origin = vertices[vIndex];
						vertices[vIndex]->leaving = tempHalfEdge;
						pList = (HalfEdgeList*)vertices[vIndex]->auxData;
						HalfEdgeList::addToList(pList, tempHalfEdge);
						vertices[vIndex]->auxData = pList;
					}
					tempHalfEdge->next = tempFace->edge;
				}

				HalfEdgeList* listWalker;
				DCELHalfEdge* eWalker = m->firstHalfEdge();
				DCELHalfEdge* newTwin;
				int numFound;
				while (isManifold && (eWalker != NULL)) {
					if (eWalker->twin == NULL) { // Haven't matched this half-edge yet.
						pList = (HalfEdgeList*)eWalker->next->origin->auxData;
						listWalker = pList;
						numFound = 0;
						newTwin = NULL;
						while (listWalker) {
							if (listWalker->edge->next->origin == eWalker->origin) {
								newTwin = listWalker->edge;
								numFound++;
							}
							listWalker = listWalker->next;
						}
						if (numFound == 0) { // Must be a boundary edge
							newTwin = new DCELHalfEdge();
							m->insert(newTwin);
							newTwin->twin = eWalker;
							eWalker->twin = newTwin;
							newTwin->face = m->getInfiniteFace();
							newTwin->origin = eWalker->next->origin;
						} else if (numFound == 1) {
							HalfEdgeList::removeFromList(pList, newTwin);
							eWalker->next->origin->auxData = pList;
							eWalker->twin = newTwin;
							newTwin->twin = eWalker;
						} else { // Two or more edges claim to originate in this list and pass through our node. This is bad
							isManifold = false;
						}
						pList = (HalfEdgeList*)eWalker->origin->auxData;
						HalfEdgeList::removeFromList(pList, eWalker);
						eWalker->origin->auxData = pList;
					}
					m->advance(eWalker);
				}

				// Even if we've decided the mesh is non-manifold, we need to clean up the auxData pointers before we clear the mesh
				// (If the mesh is manifold, this is easy. If not, we almost definitely have auxData pointers to clear)
				for (DCELVertex* vWalker = m->firstVertex(); vWalker; m->advance(vWalker)) {
					if (vWalker->auxData) { // No linked lists should be left
						isManifold = false; // If one is found, this was a bad mesh
						pList = (HalfEdgeList*)vWalker->auxData; // Clean up to prevent memory leaks from bad meshes
						HalfEdgeList::deleteList(pList);
						vWalker->auxData = NULL;
					}
				}

				if (isManifold) { // No point in doing this if we've decided it's non-manifold already, and could fail
					DCELHalfEdge* previous = NULL;
					for (eWalker = m->firstHalfEdge(); eWalker; m->advance(eWalker)) {
						if (eWalker->face == m->getInfiniteFace()) {
							previous = eWalker->twin->next->twin;
							while (previous->next) { // Note, this is dangerous if the file is very ill-formed. This could be an infinite loop
								previous = previous->next->twin;
							}
							previous->next = eWalker;
						}
					}
				}
				
				delete[] vertices;

				if (isManifold) {
					rval = true;
					m->updateAll();
				} else {
					m->clear();
				}
			}
		}

		ifs.close();
	}

	return rval;
}
