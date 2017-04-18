// DCELTools.h: interface for the DCELTools class.
//
//////////////////////////////////////////////////////////////////////

#if !defined(AFX_DCELTOOLS_H__70D80E10_60CB_4206_9FC0_8F32881E97ED__INCLUDED_)
#define AFX_DCELTOOLS_H__70D80E10_60CB_4206_9FC0_8F32881E97ED__INCLUDED_

#if _MSC_VER > 1000
#pragma once
#endif // _MSC_VER > 1000

#include "DCELMesh.h"

/*
 * DCELTools class. Part of an example DCEL implementation
 * Webpage: http://www.holmes3d.net/graphics/dcel/
 * Author: Ryan Holmes
 * E-mail: ryan <at> holmes3d <dot> net
 * Usage: Use freely. Please cite the website as the source if you
 * use it substantially unchanged. Please leave this documentation
 * in the code.
 *
 * Static function class. Demonstrates some operations on a DCEL,
 * including validation, construction, and storing to a simple file format.
 */

class DCELTools  
{
private:
	DCELTools();
	~DCELTools();

public:
	static bool isConsistent(DCELHalfEdge* e);
	static bool isConsistent(DCELFace* f);
	static bool isConsistent(DCELVertex* v);
	static bool isConsistent(DCELMesh* m);

	static bool loadFromOFF(LPCTSTR filename, DCELMesh* m);
	static bool storeToOFF(LPCTSTR filename, DCELMesh* m);
};

#endif // !defined(AFX_DCELTOOLS_H__70D80E10_60CB_4206_9FC0_8F32881E97ED__INCLUDED_)
