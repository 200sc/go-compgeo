// DCELHalfEdge.h: interface for the DCELHalfEdge class.
//
//////////////////////////////////////////////////////////////////////

#if !defined(AFX_DCELHALFEDGE_H__A8186B0F_19D5_48EF_BD30_EB266F9A8215__INCLUDED_)
#define AFX_DCELHALFEDGE_H__A8186B0F_19D5_48EF_BD30_EB266F9A8215__INCLUDED_

#if _MSC_VER > 1000
#pragma once
#endif // _MSC_VER > 1000

/*
 * DCELHalfEdge class. Part of an example DCEL implementation
 * Webpage: http://www.holmes3d.net/graphics/dcel/
 * Author: Ryan Holmes
 * E-mail: ryan <at> holmes3d <dot> net
 * Usage: Use freely. Please cite the website as the source if you
 * use it substantially unchanged. Please leave this documentation
 * in the code.
 */

class DCELFace;
class DCELVertex;

class DCELHalfEdge  
{
public:
	DCELHalfEdge();
	~DCELHalfEdge();

	enum DCEL_EDGE_BITS {
		// Set during any processing loop to indicate that the half-edge has
		// been processed. (e.g. when running through all half-edges and
		// rendering edges, set on the twin to prevent re-rendering an edge)
		// This is volatile... that is, it should not be assumed to be set or
		// unset at the beginning of a function, and need not be left in a
		// meaningful state at the end of a function
		DCEL_EDGE_PROCESSED_BIT	= 1,
		// Set to indicate that this half-edge or its twin faces the infinite face
		DCEL_EDGE_BOUNDARY_BIT = 2,
		// Set to indicate that this half-edge is selected.
		DCEL_EDGE_SELECTED_BIT = 4,
		// Set to indicate that this half-edge is marked.
		DCEL_EDGE_MARKED_BIT = 8,
		// Set to indicate that this half-edge and its twin comprise a silhouette edge.
		DCEL_EDGE_SILHOUETTE_BIT = 16
	};

	DCELHalfEdge* twin;
	DCELHalfEdge* next;
	DCELFace* face;
	DCELVertex* origin;
	void* auxData;
	// Optional member to store some state during operations. Not maintained internally,
	// except by user-initiated calls on the full mesh.
	unsigned int displayBits;

	DCELHalfEdge* getPrev();
	// Helper functions for manipulating the displayBits property
	bool isMaskSet(unsigned int mask) const { return ((displayBits & mask) == mask); };
	void setMask(unsigned int mask, bool value) { displayBits = (value ? (displayBits | mask) : (displayBits & ~mask)); };

	friend class DCELMesh;
protected:
	DCELHalfEdge* globalNext;
	DCELHalfEdge* globalPrev;
};

#endif // !defined(AFX_DCELHALFEDGE_H__A8186B0F_19D5_48EF_BD30_EB266F9A8215__INCLUDED_)
