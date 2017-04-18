#pragma once

#include "DCELHalfEdge.h"

/*
 * HalfEdgeList class. Part of an example DCEL implementation
 * Webpage: http://www.holmes3d.net/graphics/dcel/
 * Author: Ryan Holmes
 * E-mail: ryan <at> holmes3d <dot> net
 * Usage: Use freely. Please cite the website as the source if you
 * use it substantially unchanged. Please leave this documentation
 * in the code.
 *
 * Simple structure to manipulate a singly linked list of DCELHalfEdge*s
 * Used as temporary structure hung from auxData pointers during a
 * mesh-building operation.
 */


class HalfEdgeList
{
public:
	HalfEdgeList(void);
	~HalfEdgeList(void);

	DCELHalfEdge* edge;
	HalfEdgeList* next;

	static void addToList(HalfEdgeList* &head, DCELHalfEdge* newEdge);
	static void addToList(HalfEdgeList* &head, HalfEdgeList* newItem);
	static void deleteList(HalfEdgeList* &head);
	static int getListLength(HalfEdgeList* head);
	static bool removeFromList(HalfEdgeList* &head, DCELHalfEdge* edge);

};
